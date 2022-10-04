package email

import "context"

type Email interface {
	SendEmail(ctx context.Context, template, destination, subject string, data interface{}) error
}
