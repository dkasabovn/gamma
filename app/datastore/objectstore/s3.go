package objectstore

import (
	"bytes"
	"context"
	"gamma/app/external/aws"
	"gamma/app/system"
	"gamma/app/system/log"
	"gamma/app/system/util"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type storageS3 struct {
	handle *s3.Client
}

func newS3() *storageS3 {
	client := s3.NewFromConfig(aws.NewConfig())
	return &storageS3{
		handle: client,
	}
}

func (s *storageS3) Get(ctx context.Context, key string) (*Object, error) {
	obj, err := s.handle.GetObject(ctx, &s3.GetObjectInput{
		Bucket: util.StringRef(system.GetConfig().BucketName),
		Key:    util.StringRef(key),
	})

	if err != nil {
		log.Infof("Could not get object: %v; %v", key, err)
		return nil, err
	}

	data, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	return &Object{
		Data: data,
		Key:  key,
	}, nil
}

func (s *storageS3) Put(ctx context.Context, key string, obj *Object) (string, error) {
	_, err := s.handle.PutObject(ctx, &s3.PutObjectInput{
		Bucket: util.StringRef(system.GetConfig().BucketName),
		Key:    util.StringRef(key),
		Body:   bytes.NewReader(obj.Data),
	})
	return key, err
}
