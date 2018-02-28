package s3store

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/margic/bazel-s3-cache/store"
	"github.com/prometheus/client_golang/prometheus"
)

var _ store.Storer = (*S3Store)(nil)

// S3Store is an implementation of Storer a backend store for the bazel remote cache
type S3Store struct {
	bucket         string
	dl             *s3manager.Downloader
	ul             *s3manager.Uploader
	cacheDurations *prometheus.SummaryVec
	errCounter     prometheus.Counter
}

// NewS3Store returns a new s3 store that implements Storer
func NewS3Store(region, bucket, key, secret, token string) (*S3Store, error) {
	var sess *session.Session
	if len(key) > 0 || len(secret) > 0 || len(token) > 0 {
		s, err := session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(key, secret, token),
		},
		)
		if err != nil {
			return nil, err
		}
		sess = s
	} else {
		s, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
		},
		)
		if err != nil {
			return nil, err
		}
		sess = s
	}

	down := s3manager.NewDownloader(sess)
	up := s3manager.NewUploader(sess)

	cacheDurations := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "cache_durations_seconds",
			Namespace:  "bazel",
			Subsystem:  "cache",
			Help:       "Cache latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path"},
	)

	errorCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:      "cache_error_count",
			Help:      "Count of Bazel cache errors",
			Namespace: "bazel",
			Subsystem: "cache",
		},
	)
	prometheus.MustRegister(cacheDurations)
	prometheus.MustRegister(errorCounter)

	s3Store := &S3Store{
		dl:             down,
		ul:             up,
		bucket:         bucket,
		cacheDurations: cacheDurations,
		errCounter:     errorCounter,
	}

	log.Printf("s3store bucket: %s", bucket)
	return s3Store, nil
}

// Get return object from the store
func (s *S3Store) Get(ctx context.Context, path string, sha string, w io.WriterAt) (int64, error) {
	if !validPath(path) {
		s.errCounter.Inc()
		return 0, fmt.Errorf("invalid path %s", path)
	}
	start := time.Now()
	defer s.cacheDurations.WithLabelValues("GET", path).Observe(float64(time.Since(start).Seconds()))
	in := &s3.GetObjectInput{
		Key:    aws.String(path + "/" + sha),
		Bucket: aws.String(s.bucket),
	}
	size, err := s.dl.DownloadWithContext(ctx, w, in)
	if err != nil {
		s.countError(err)
	}
	return size, err
}

// Put put oject into the store
func (s *S3Store) Put(ctx context.Context, path string, sha string, r io.Reader) error {
	if !validPath(path) {
		s.errCounter.Inc()
		return fmt.Errorf("invalid path %s", path)
	}
	start := time.Now()
	defer s.cacheDurations.WithLabelValues("PUT", path).Observe(float64(time.Since(start).Seconds()))
	in := &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Body:   r,
		Key:    aws.String(path + "/" + sha),
	}
	_, err := s.ul.UploadWithContext(ctx, in)
	if err != nil {
		s.countError(err)
	}
	return err
}

func (s *S3Store) countError(err error) {
	var s3Err, ok = err.(s3.RequestFailure)
	if ok {
		if s3Err.StatusCode() == 404 {
			// not an error
			return
		}
	}
	s.errCounter.Inc()
}

func validPath(path string) bool {
	if path == "ac" || path == "cas" {
		return true
	}
	return false
}
