package objectstore

import (
	"context"

	"cloud.google.com/go/storage"
)

type cloudStorageClient struct {
	client *storage.Client
}

func NewCloudStorage(ctx context.Context) *cloudStorageClient {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return &cloudStorageClient{
		client: client,
	}
}

func (gc *cloudStorageClient) PutObject(ctx context.Context, key string, data []byte) error {
	objHandle := gc.client.Bucket("").Object(key).NewWriter(ctx)
	objHandle.Write(data)
	if err := objHandle.Close(); err != nil {

	}
	return nil
}

func (gc *cloudStorageClient) GetObjectUrl(ctx context.Context, key string) (string, error) {
	url, err := storage.SignedURL("", key, &storage.SignedURLOptions{
		Method: "GET",
	})

	if err != nil {
		return "", err
	}

	return url, nil
}
