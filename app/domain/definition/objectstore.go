package definition

import "context"

type Store interface {
	PutObject(ctx context.Context, key string, data []byte) error
	GetObjectUrl(ctx context.Context, key string)
}
