package dispatch

import (
	"fmt"
	"uno/pkg/setting"
	"uno/service/provider"
	"uno/service/subscriber"
	"uno/service/template"

	"github.com/gin-gonic/gin"
)

// WorkflowCreate is the entry for workflow creation.
func WorkflowCreate(ctx *gin.Context) {
	// . Read topic.

	// . TODO Read the subscribers.
	subscribers := make([]*subscriber.Entry, 0)
	subscribers = append(subscribers, &subscriber.Entry{
		Email: "hi@noahyao.me",
	}, &subscriber.Entry{
		// UserID => 8000000000001042
		Email: "xlyiin227@163.com",
	})

	// . Append additional info for subscriber.
	for _, subscriber := range subscribers {
		// TODO call passport
		fmt.Println(subscriber.Email)
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

	// . Send the message.
	for _, subscriber := range subscribers {
		for _, provider := range providers {
			tpl, err := template.Get("marketing-20231117", map[string]string{"hello": "world"})
			if err != nil {
				fmt.Println(err)
				continue
			}

			providerResponse, err := provider.Send(ctx, subscriber, tpl)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(subscriber.Email, subscriber.UserID, "message_id", providerResponse.MessageID)
		}
	}

	// . Mark the subscriber result.
}
