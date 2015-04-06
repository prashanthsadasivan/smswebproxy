package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"smswebproxy/app/models"
	"strings"
)

func sendMessage(message models.SMSMessage, reg string) {
	fmt.Printf("sendingMessage: %s\n", message.Message)

	client := http.Client{}

	payload, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	datastring := string(payload[:])
	fmt.Printf("datastring: %s\n", datastring)
	fmt.Printf("regid: %s\n", reg)
	fmt.Printf("gcm: %s\n", os.Getenv("GCM_AUTH_KEY"))

	request, err := http.NewRequest("POST", "https://android.googleapis.com/gcm/send", strings.NewReader("{\"registration_ids\":[\""+reg+"\"], \"data\" : "+datastring+"}"))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "key="+os.Getenv("GCM_AUTH_KEY"))
	resp, err2 := client.Do(request)
	if err2 != nil {
		panic(err2)
	}
	if resp != nil {
		fmt.Printf("responseCode: %d\n", resp.StatusCode)
	}
}

func (c App) Websock(num string, ws *websocket.Conn) revel.Result {
	// In order to select between websocket messages and subscription events, we
	// need to stuff websocket events into a channel.
	fmt.Printf("ws num: %s\n", num)
	appconnection := c.GorcAppConnection(num)
	messagesToSend := make(chan models.SMSMessage)
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				fmt.Printf("err: %s\n", err.Error())
				close(messagesToSend)
				return
			}
			fmt.Printf("got message from ws: %s\n", msg)
			if msg == "ping" {
				if websocket.JSON.Send(ws, "pong") != nil {
					// They disconnected.
					return
				}
			} else {
				var sms models.SMSMessage
				json.Unmarshal([]byte(msg), &sms)
				messagesToSend <- sms
			}
		}
	}()

	// Now listen for new events from either the websocket or the chatroom.
	for {
		select {
		case event := <-appconnection.Received:
			fmt.Printf("ws, sent on websocket\n")
			if websocket.JSON.Send(ws, &event) != nil {
				// They disconnected.
				return nil
			}
		case msg, ok := <-messagesToSend:
			// If the channel is closed, they disconnected.
			fmt.Printf("ws sent a message")
			if !ok {
				return nil
			}
			// Otherwise, say something.
			sendMessage(msg, appconnection.RegId)
		}
	}
	return nil
}
