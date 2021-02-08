package saver

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/crawler/src/extractor"
	"log"
	"time"
)

func UploadScrapeToBucket(bucketName string) {

	currentTime := time.Now()
	timeFormat := currentTime.Format("2-Jan-2006-15:04")

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Panic(err)
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object("write_p/"+timeFormat)

	w := obj.NewWriter(ctx)
	if _, err := fmt.Fprintf(w, "%s", extractor.Extract()); err != nil {
		log.Panic(err)
	}
	if err := w.Close(); err != nil {
		log.Panic(err)
	}

}

