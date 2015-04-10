package websocket_routing

import (
	"golang.org/x/net/websocket"
	"smswebproxy/app/appconnections"
	"smswebproxy/app/models"
)

func handlePing(sms models.SMSMessage, appconnection *appconnections.AppConnection, ws *websocket.Conn) bool {
	return websocket.JSON.Send(ws, "pong") != nil
}
