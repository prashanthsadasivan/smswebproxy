package websocket_routing

import (
	"golang.org/x/net/websocket"
	"smswebproxy/app/appconnections"
	"smswebproxy/app/models"
)

type WebsocketMessageHandler func(sms models.SMSMessage, appconnection *appconnections.AppConnection, ws *websocket.Conn) bool

var Router map[string]WebsocketMessageHandler

func init() {
	Router = make(map[string]WebsocketMessageHandler)
	Router["ping"] = handlePing
	Router["SMS/Send"] = sendSmsMessage
}
