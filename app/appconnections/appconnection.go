package appconnections

import (
	"smswebproxy/app/models"
)

var (
	conduits map[string]*AppConnection
)

func init() {
	conduits = make(map[string]*AppConnection)
}

type AppConnection struct {
	Received chan models.SMSMessage
	RegId    string
}

func GetAppConnection(number string) *AppConnection {
	return conduits[number]
}

func AddAppConnection(num string, conduit *AppConnection) {
	conduits[num] = conduit
}

func New(regid, num string) *AppConnection {
	conduit := new(AppConnection)
	conduit.RegId = regid
	conduit.Received = make(chan models.SMSMessage)
	conduits[num] = conduit
	return conduit
}
