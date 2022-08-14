package email

import (
	"context"
	"sync"

	"gamma/app/external/aws"
	"gamma/app/system/util"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

var (
	emailOnce sync.Once
	emailInst *Email
)

type Email struct {
	sesHandle *ses.Client
}

func GetEmail() *Email {
	emailOnce.Do(func() {
		emailInst = &Email{
			sesHandle: ses.NewFromConfig(aws.NewConfig()),
		}
	})
	return emailInst
}

func (e *Email) SendEmail(ctx context.Context, destination string) {
	e.sesHandle.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{destination},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    util.StringRef("Subject"),
				Charset: util.StringRef("UTF-8"),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data:    util.StringRef("TEMPLATE HERE"),
					Charset: util.StringRef("UTF-8"),
				},
				Text: &types.Content{
					Data:    util.StringRef("TEMPLATE HERE"),
					Charset: util.StringRef("UTF-8"),
				},
			},
		},
		Source: new(string),
	})
}
