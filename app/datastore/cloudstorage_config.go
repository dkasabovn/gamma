package datastore

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var BUCKET_NAME string
var csCon *storage.Client
var CSSingleton sync.Once

func StorageInstance() *storage.Client {
	CSSingleton.Do(func() {
		client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(os.Getenv("CLOUD_CREDENTIALS_FILE")))
		if err != nil {
			panic("Could not open file")
		}
		BUCKET_NAME = os.Getenv("BUCKET_NAME")
		csCon = client
	})
	return csCon
}
