package dispatch

import (
	"fmt"
	"time"
	"uno/pkg/database"
	"uno/pkg/setting"
	"uno/service/message"
	"uno/service/provider"
	"uno/service/template"
	"uno/service/topic"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Workflow struct {
	ID         string `json:"id,omitempty"`
	TopicID    string `json:"topic_id,omitempty"`
	TemplateID string `json:"template_id,omitempty"`
}

// WorkflowCreate is the entry for workflow creation.
func WorkflowCreate(ctx *gin.Context) {
	workflow := &Workflow{}
	if ctx.BindJSON(workflow) != nil {
		ctx.JSON(400, gin.H{"message": "invalid request"})
		return
	}

	workflow.ID = uuid.New().String()

	// . Read topic.
	topic, _ := topic.Get(workflow.TopicID)
	if topic == nil {
		ctx.JSON(400, gin.H{"message": "topic not found", "topic_id": workflow.TopicID})
		return
	}

	if topic.Status != 1 {
		ctx.JSON(400, gin.H{"message": "topic is not active"})
		return
	}

	// . Initialize the template.

	// . Initialize the providers.
	providers := make([]provider.Base, 0)
	providers = append(providers, provider.NewAWSEmail(
		setting.AppInstance.Provider.EmailSES.Access,
		setting.AppInstance.Provider.EmailSES.Secret,
		setting.AppInstance.Provider.EmailSES.Region))

	for _, provider := range providers {
		provider.SetOption(ctx, nil)
	}

	sentCount := 0

	// . Send the message.
	for _, subscriber := range topic.Subscribers {
		for _, provider := range providers {
			tpl, err := template.Get(workflow.TemplateID, map[string]string{"hello": "world"})
			if err != nil {
				fmt.Println(err)
				continue
			}

			digest := provider.Digest(subscriber, tpl)
			msg := &message.Entry{}
			database.GetWriteDB().Select("id").Where("digest = ?", digest).First(msg)
			if msg.ID != "" {
				fmt.Println("message already sent", msg.ID, msg.Digest)
				continue
			}

			providerResponse, err := provider.Send(ctx, subscriber, tpl)
			if err != nil {
				fmt.Printf("failed to send message: %v\n", err)
				continue
			}

			sentCount++

			msg.ID = uuid.New().String()
			msg.WorkflowID = workflow.ID
			msg.UserID = subscriber.UserID
			msg.Digest = providerResponse.Digest
			msg.Channel = "aws_email_ses"
			msg.ChannelMessageID = providerResponse.MessageID
			msg.ChannelIdentifier = subscriber.Email
			msg.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
			msg.Status = 1
			if database.GetWriteDB().Create(msg).Error != nil {
				fmt.Printf("failed to create message: %v\n", msg)
			}

			fmt.Println(subscriber.Email, subscriber.UserID, "message_id", providerResponse.MessageID, "digest", providerResponse.Digest)
		}
	}

	// . Mark the subscriber result.

	ctx.JSON(200, gin.H{"message": "ok", "workflow": workflow, "sent_cnt": sentCount})
}
