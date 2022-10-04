package email

import (
	"context"
	"sync"

	"gamma/app/external/aws"
	"gamma/app/system/log"
	"gamma/app/system/util"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

var (
	emailOnce sync.Once
	emailInst *email
)

type email struct {
	sesHandle *ses.Client
}

func GetEmail() Email {
	emailOnce.Do(func() {
		emailInst = &email{
			sesHandle: ses.NewFromConfig(aws.NewConfig()),
		}
	})
	return emailInst
}

func (e *email) SendEmail(ctx context.Context, template, destination, subject string, data interface{}) error {
	htmlTmpl, err := render(template, html, data)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	txtTmpl, err := render(template, text, data)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	if _, err := e.sesHandle.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{destination},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    util.StringRef(subject),
				Charset: util.StringRef("UTF-8"),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data:    util.StringRef(htmlTmpl),
					Charset: util.StringRef("UTF-8"),
				},
				Text: &types.Content{
					Data:    util.StringRef(txtTmpl),
					Charset: util.StringRef("UTF-8"),
				},
			},
		},
		Source: new(string),
	}); err != nil {
		log.Errorf("%v", err)
		return err
	}

	return nil
}
