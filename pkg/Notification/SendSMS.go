package Notification

import (
	"encoding/json"
	"fmt"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go-ecommerce-app/Config"
)

type NotificatinClient interface {
	SendSMS(phone string, message string) error
}

type notificatinClient struct {
	config Config.AppConfig
}

// twilio
func (c notificatinClient) SendSMS(phone string, message string) error {
	accountSid := c.config.AccountSid
	authToken := c.config.AuthToken

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(c.config.Fromphone) //from twilio
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
	return nil
}

func NewNotificationClient(config Config.AppConfig) NotificatinClient {
	return &notificatinClient{
		config: config,
	}
}
