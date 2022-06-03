package objectstore

import (
	"gamma/app/datastore"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/gommon/log"
)

func GenerateObjectUploadUrl(objectPath string) (string, error) {
	url, err := datastore.StorageInstance().Bucket(datastore.BUCKET_NAME).GenerateSignedPostPolicyV4(objectPath, &storage.PostPolicyV4Options{
		Expires: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		log.Errorf("could not generate pre-signed url: %v", err)
		return "", err
	}
	return url.URL, nil
}
