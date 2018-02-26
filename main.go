package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/margic/bazel-s3-cache/api"
	"github.com/margic/bazel-s3-cache/s3store"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	flag.StringP("listen", "l", ":8080", "Set the local address to listen on")
	viper.BindEnv("listen", "LISTEN_ADDR")
	viper.BindPFlag("listen", flag.Lookup("listen"))
	flag.StringP("bucket", "b", "", "Set the s3 bucket name")
	viper.BindEnv("bucket", "AWS_S3_BUCKET")
	viper.BindPFlag("bucket", flag.Lookup("bucket"))
	flag.StringP("region", "r", "us-east-1", "Set region for s3 bucket")
	viper.BindEnv("region", "AWS_DEFAULT_REGION")
	viper.BindPFlag("region", flag.Lookup("region"))

	flag.StringP("key", "k", "", "Set aws access key")
	viper.BindEnv("key", "AWS_ACCESS_KEY_ID")
	viper.BindPFlag("key", flag.Lookup("key"))
	flag.StringP("secret", "s", "", "Set aws secret key")
	viper.BindEnv("secret", "AWS_SECRET_ACCESS_KEY")
	viper.BindPFlag("secret", flag.Lookup("secret"))
	flag.StringP("token", "t", "", "Set aws sessiont token")
	viper.BindEnv("token", "AWS_SESSION_TOKEN")
	viper.BindPFlag("token", flag.Lookup("token"))
	flag.Parse()

	if !viper.IsSet("bucket") || len(viper.GetString("bucket")) == 0 {
		log.Fatal("s3 bucket for cache not set")
	}
	addr := viper.GetString("listen")
	region := viper.GetString("region")
	bucket := viper.GetString("bucket")
	key := viper.GetString("key")
	secret := viper.GetString("secret")
	token := viper.GetString("token")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// make a backend store
	stor, err := s3store.NewS3Store(region, bucket, key, secret, token)
	if err != nil {
		log.Printf("error creating new s3 store: %s", err.Error())
	}

	// pass store to api server
	apiServer, err := api.NewServer(stor)
	if err != nil {
		log.Printf("error creating server: %s", err.Error())
	}

	// transport juicyness
	httpServer := http.Server{
		Addr:    addr,
		Handler: apiServer.Handler,
	}

	log.Printf("starting cache http server on %s", addr)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// wait for interrupt
	<-stop
	// shutdown httpServer
	ctx := context.Background()
	httpServer.Shutdown(ctx)

	log.Println("bazel cache server stopped")
}
