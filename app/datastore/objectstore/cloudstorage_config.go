package objectstore

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var csCon *storage.Client
var storageSingleton sync.Once

func storageInstance() *storage.Client {
	storageSingleton.Do(func() {
		client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(os.Getenv("CLOUD_CREDENTIALS_FILE")))
		if err != nil {
			panic("Could not open file")
		}
		csCon = client
	})
	return csCon
}
