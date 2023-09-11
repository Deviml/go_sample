package ses

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	sesv2 "github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/go-kit/kit/log"
)

const (
	CharSet = "UTF-8"
)

type Client struct {
	sender string
	ses    *ses.Client
	sesv2  *sesv2.SESV2
	logger log.Logger
}

func NewClient(sender string, ses *ses.Client, sesv2 *sesv2.SESV2, logger log.Logger) *Client {
	return &Client{sender: sender, ses: ses, sesv2: sesv2, logger: logger}
}

func (c Client) SendEmail(ctx context.Context, to []string, body string, subject string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: to,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(body),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(c.sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	request := c.ses.SendEmailRequest(input)
	_, err := request.Send(ctx)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				c.logger.Log(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				c.logger.Log(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				c.logger.Log(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				c.logger.Log(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			c.logger.Log(err.Error())
		}

		return err
	}
	return nil
}

func (c Client) SendEmailInBulkV2(ctx context.Context, to []string, name string, category string, amount int, city string, zipcode string, specialRequest string, state string, FrontURL string, ID uint) error {

	entries := make([]*sesv2.BulkEmailEntry, len(to))
	for i, recipient := range to {
		entry := &sesv2.BulkEmailEntry{
			Destination: &sesv2.Destination{
				ToAddresses: aws.StringSlice([]string{recipient}),
			},
			ReplacementEmailContent: &sesv2.ReplacementEmailContent{
				ReplacementTemplate: &sesv2.ReplacementTemplate{
					ReplacementTemplateData: aws.String(fmt.Sprintf("{ \"Name\":\"%s\", \"Category\": \"%s\", \"Amount\":\"%d\", \"City\":\"%s\",  \"Zipcode\":\"%s\",  \"SpecialRequest\":\"%s\",  \"State\":\"%s\",  \"Link\":\"%s\",  \"ID\":\"%d\", }", name, category, amount, city, zipcode, specialRequest, state, FrontURL, ID)),
				},
			},
		}
		entries[i] = entry
	}

	input := &sesv2.SendBulkEmailInput{
		FromEmailAddress: aws.String(c.sender),
		BulkEmailEntries: entries,
		DefaultContent: &sesv2.BulkEmailContent{
			Template: &sesv2.Template{
				TemplateName: aws.String("emailNotificationsV2"),
				TemplateData: aws.String("{ \"Name\":\"User\", \"Category\": \"Aerial Work Platform Equipments\", \"Amount\":\"00\", \"City\":\"Miami\",  \"Zipcode\":\"12345\",  \"SpecialRequest\":\"Test Request\",  \"State\":\"Florida\",  \"Link\":\"www.equiphunter.com\",  \"ID\":\"1\", }"),
			},
		},
	}

	output, err := c.sesv2.SendBulkEmailWithContext(ctx, input)
	if err != nil {
		c.logger.Log(err.Error())
		return err
	}

	failedCount := 0
	for _, entry := range output.BulkEmailEntryResults {
		if entry.Status != nil && *entry.Status == sesv2.BulkEmailStatusFailed {
			failedCount++
			c.logger.Log("Failed email: " + *entry.Error)
		}
	}
	if failedCount > 0 {
		errMsg := fmt.Sprintf("%d emails failed to send", failedCount)
		return errors.New(errMsg)
	}
	return nil
}

func (c Client) SendForgotPassword(ctx context.Context, to string, password string) error {

	emailAddress := make([]*string, len(to))

	input := &sesv2.SendEmailInput{
		Content: &sesv2.EmailContent{
			Template: &sesv2.Template{
				TemplateName: aws.String("ForgotPasswordV1"),
				TemplateData: aws.String(fmt.Sprintf("{ \"Password\":\"%s\" }", password)),
			},
		},
		Destination: &sesv2.Destination{
			ToAddresses: emailAddress,
		},
		FromEmailAddress: aws.String(c.sender),
	}

	_, err := c.sesv2.SendEmailWithContext(ctx, input)
	if err != nil {
		c.logger.Log(err.Error())
		return err
	}

	return nil
}