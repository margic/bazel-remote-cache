package s3store

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/margic/bazel-s3-cache/store"
)

var _ store.Storer = (*S3Store)(nil)

// S3Store is an implementation of Storer a backend store for the bazel remote cache
type S3Store struct {
	bucket string
	dl     *s3manager.Downloader
	ul     *s3manager.Uploader
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
	s3Store := &S3Store{
		dl:     down,
		ul:     up,
		bucket: bucket,
	}
	log.Printf("s3store bucket: %s", bucket)
	return s3Store, nil
}

// Get return object from the store
func (s *S3Store) Get(ctx context.Context, path string, sha string, w io.WriterAt) (int64, error) {
	if !validPath(path) {
		return 0, fmt.Errorf("invalid path %s", path)
	}
	in := &s3.GetObjectInput{
		Key:    aws.String(path + "/" + sha),
		Bucket: aws.String(s.bucket),
	}
	log.Printf("%s", *in.Key)
	return s.dl.DownloadWithContext(ctx, w, in)
}

// Put put oject into the store
func (s *S3Store) Put(ctx context.Context, path string, sha string, r io.Reader) error {
	if !validPath(path) {
		return fmt.Errorf("invalid path %s", path)
	}
	in := &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Body:   r,
		Key:    aws.String(path + "/" + sha),
	}
	log.Printf("%s", *in.Key)
	_, err := s.ul.UploadWithContext(ctx, in)
	if err != nil {
		return err
	}
	return nil
}

func validPath(path string) bool {
	if path == "ac" || path == "cas" {
		return true
	}
	return false
}
