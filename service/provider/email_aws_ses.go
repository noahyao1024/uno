package provider

import (
	"fmt"
	"uno/service/subscriber"
	"uno/service/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
)

type SESOfAWS struct {
	Option *Option

	AK, SK, Region string
}

func NewAWSEmail(ak, sk, region string) *SESOfAWS {
	return &SESOfAWS{
		Option: &Option{},
		AK:     ak,
		SK:     sk,
		Region: region,
	}
}

func (i *SESOfAWS) SetOption(ctx *gin.Context, o *Option) error {
	return nil
}

func (i *SESOfAWS) Send(ctx *gin.Context, subscriber *subscriber.Entry, template *template.Entry) (*Response, error) {
	if len(subscriber.Email) == 0 {
		return nil, fmt.Errorf("subscriber email is empty")
	}

	if len(template.Subject) == 0 || len(template.Content) == 0 {
		return nil, fmt.Errorf("template subject or content is empty")
	}

	awsSession, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(i.AK, i.SK, ""),
		Region:      aws.String(i.Region),
	})

	if err != nil {
		return nil, err
	}

	sesClient := ses.New(awsSession)

	inputConfig := &ses.SendEmailInput{
		Source: aws.String(template.Sender), // 发件箱
		Destination: &ses.Destination{
			BccAddresses: nil,
			CcAddresses:  []*string{&subscriber.Email},
			ToAddresses:  []*string{&subscriber.Email},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(template.Content),
				},
				Text: nil,
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String(template.Subject),
			},
		},
	}

	sendEmailOutput, err := sesClient.SendEmail(inputConfig)

	if err != nil {
		return nil, err
	}

	if sendEmailOutput.MessageId == nil {
		return nil, fmt.Errorf("send email failed, message id is empty")
	}

	return &Response{MessageID: *sendEmailOutput.MessageId}, nil
}
