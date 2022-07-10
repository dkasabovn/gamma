package objectstore

import (
	"context"
	"gamma/app/datastore"
	"gamma/app/system"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/gommon/log"
)

type cloudStore struct {
	bucket *storage.BucketHandle
}

func GenerateObjectUploadUrl(objectPath string) (string, error) {
	url, err := datastore.StorageInstance().Bucket(datastore.BUCKET_NAME).GenerateSignedPostPolicyV4(objectPath, &storage.PostPolicyV4Options{
		Expires:    time.Now().Add(15 * time.Minute),
		Conditions: []storage.PostPolicyV4Condition{},
	})
	if err != nil {
		log.Errorf("could not generate pre-signed url: %v", err)
		return "", err
	}
	return url.URL, nil
}

func newCloudStore() Storage {
	return &cloudStore{
		bucket: datastore.StorageInstance().Bucket(system.BUCKET_NAME),
	}
}

func (c *cloudStore) Put(ctx context.Context, key string, put Object) (string, error) {
	w := c.bucket.Object(key).NewWriter(ctx)
	defer w.Close()
	if _, err := w.Write(put.Data); err != nil {
		return "", err
	}
	return w.Name, nil
}

func (c *cloudStore) Get(ctx context.Context, key string) (Object, error) {
	r, err := c.bucket.Object(key).NewReader(ctx)
	if err != nil {
		return Object{}, err
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Object{}, err
	}

	return Object{
		Data: data,
		Key:  key,
	}, nil
}
