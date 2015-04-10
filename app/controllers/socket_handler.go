package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"smswebproxy/app/controllers/websocket_routing"
)

func (c App) Websock(num string, ws *websocket.Conn) revel.Result {
	// In order to select between websocket messages and subscription events, we
	// need to stuff websocket events into a channel.
	fmt.Printf("ws num: %s\n", num)
	appconnection := c.GorcAppConnection(num)
	messagesToSend := appconnection.Start(ws)
	fmt.Printf("socket appconn %+v\n", appconnection)

	// Now listen for new events from either the websocket or the chatroom.
	for {
		select {
		case event := <-appconnection.Received:
			fmt.Printf("ws, sent on websocket: %+v\n", event)
			if websocket.JSON.Send(ws, &event) != nil {
				// They disconnected.
				fmt.Printf("here")
				return nil
			}
		case msg, ok := <-messagesToSend:
			// If the channel is closed, they disconnected.
			fmt.Printf("ws sent a message")
			if !ok {
				return nil
			}
			websocket_routing.Router[msg.MessageType](msg, appconnection, ws)
		}
	}
	return nil
}
