package ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"notification-service/internal/core/domain"
)

// Sender implementa output.EmailSender usando Amazon SES v2.
type Sender struct {
	client *sesv2.Client
	from   string
}

func NewSender(cfg aws.Config, fromEmail string) *Sender {
	return &Sender{
		client: sesv2.NewFromConfig(cfg),
		from:   fromEmail,
	}
}

func (s *Sender) Send(ctx context.Context, email domain.Email) error {
	_, err := s.client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(s.from),
		Destination: &types.Destination{
			ToAddresses: []string{email.To},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(email.Subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(email.HTML),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ses SendEmail: %w", err)
	}

	return nil
}
