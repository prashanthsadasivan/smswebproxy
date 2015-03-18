package models

import (
	"time"
)

type Message struct {
	Num        string
	Message    string
	MessageId  int
	CreateDate time.Time
}

type SoleUser struct {
	Number string
	GcmId  string
	UserId int
}
