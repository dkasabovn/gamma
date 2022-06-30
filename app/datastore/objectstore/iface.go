package objectstore

import "context"

type Object struct {
	Data []byte
	Key  string
}

type Storage interface {
	Put(ctx context.Context, key string, obj Object) (string, error)
	Get(ctx context.Context, key string) (Object, error)
	// Delete(ctx context.Context, key string) error
}
