package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"smswebproxy/app/models"
	"strings"
)

func (c App) Receive(receiver, num_from, messageReceived string) revel.Result {
	fmt.Printf("received text message %s %s %s", receiver, num_from, messageReceived)
	sms := models.SMSMessage{Num: strings.TrimPrefix(num_from, "+1"), Message: messageReceived}
	appconnection := c.GorcAppConnection(receiver)
	fmt.Printf("appconn %+v\n", appconnection)
	go func() {
		fmt.Printf("beforeSend")
		appconnection.Received <- sms
		fmt.Printf("received, sent on channel\n")
	}()
	c.Response.Status = 202
	return c.RenderText("")
}
