package controllers

import (
	"bytes"
	"golang.org/x/net/websocket"
	"code.google.com/p/rsc/qr"
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"net/http"
	"os"
	"smswebproxy/app/models"
	"smswebproxy/app/room"
	"strings"
	"time"
)

type App struct {
	GorpController
}

func makeConduit(regid, num string) *room.Conduit {
	conduit := new(room.Conduit)
	conduit.RegId = regid
	conduit.Received = make(chan room.SMSMessage)
	room.AddConduit(num, conduit)
	return conduit
}

func (c App) Gcm(regid, num string) revel.Result {
	fmt.Printf("regId: %s\n num:%s\n postnum: %s\n", regid, num, c.Request.PostFormValue("num"))
	makeConduit(regid, num)
	results, err := c.Txn.Select(models.SoleUser{}, "select * from SoleUser")
	if err != nil {
		panic(err)
	}
	fmt.Printf("results: %s\n", results)
	if len(results) >= 1 {
		numDeleted, err := c.Txn.Delete(results...)
		if err != nil {
			fmt.Printf("err: %s\n", err)
		}
		fmt.Printf("numDeleted: %d\n", numDeleted)
	}
	newSoleUser := &models.SoleUser{Number: num, GcmId: regid}
	insertErr := c.Txn.Insert(newSoleUser)
	if insertErr != nil {
		fmt.Printf("here err: %s\n", insertErr)
	}
	c.Response.Status = 201
	return c.RenderText("created")
}

func (c App) Send(sender, num_to, msgToTextOut string) revel.Result {
	return nil
}

func sendMessage(message room.SMSMessage, reg string) {
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

func (c App) QR(hostname string) revel.Result {
	code, err := qr.Encode(hostname, qr.H)
	if err != nil {
		panic(err)
	}
	png := code.PNG()
	modtime := time.Now()
	c.Response.Out.Header().Set("Content-Type", "image/png")
	return &revel.BinaryResult{
		Reader:  bytes.NewReader(png),
		Name:    "qr.png",
		ModTime: modtime,
	}
}

func (c App) Receive(receiver, num_from, messageReceived string) revel.Result {
	fmt.Printf("received text message")
	conduit := room.GetConduit(receiver)
	if conduit == nil {
		soleUser := &models.SoleUser{}
		err := c.Txn.SelectOne(soleUser, "select * from SoleUser")
		if err != nil {
			panic(err)
		}
		conduit = makeConduit(soleUser.GcmId, soleUser.Number)
	}
	sms := room.SMSMessage{Num: num_from, Message: messageReceived}
	conduit.Received <- sms
	fmt.Printf("received, sent on channel\n")
	c.Response.Status = 201
	return c.RenderText("")
}

func (c App) Room(num string) revel.Result {
	return c.Render(num)
}

func (c App) Home() revel.Result {
	return c.Render()
}

func (c App) Websock(num string, ws *websocket.Conn) revel.Result {
	// In order to select between websocket messages and subscription events, we
	// need to stuff websocket events into a channel.
	fmt.Printf("ws num: %s\n", num)
	conduit := room.GetConduit(num)
	if conduit == nil {
		soleUser := &models.SoleUser{}
		err := c.Txn.SelectOne(soleUser, "select * from SoleUser")
		if err != nil {
			panic(err)
		}
		conduit = makeConduit(soleUser.GcmId, soleUser.Number)
	}
	messagesToSend := make(chan room.SMSMessage)
	go func() {
    var msg string
		for {
      fmt.Printf("%+v\n", ws)
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
				var sms room.SMSMessage
				json.Unmarshal([]byte(msg), &sms)
				messagesToSend <- sms
			}
		}
	}()

	// Now listen for new events from either the websocket or the chatroom.
	for {
		select {
		case event := <-conduit.Received:
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
			sendMessage(msg, conduit.RegId)
		}
	}
	return nil
}
