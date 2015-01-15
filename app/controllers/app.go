package controllers

import (
    "github.com/revel/revel"
    "smswebproxy/app/room"
    "fmt"
    "net/http"
    "encoding/json"
    "strings"
    "code.google.com/p/go.net/websocket"
)

type App struct {
    *revel.Controller
}

func (c App) Gcm(regid, num string) revel.Result {
    fmt.Printf("regId: %s\n num:%s\n", regid, num)
    conduit := new(room.Conduit)
    conduit.RegId = regid
    conduit.Received = make(chan room.SMSMessage)
    room.AddConduit(num, conduit)
    c.Response.Status = 201
    return nil
}

func (c App) Send(sender, num_to, msgToTextOut string) revel.Result {
    return nil
}

func sendMessage(message room.SMSMessage, reg string) {
    fmt.Printf("sendingMessage: %s\n", message.Message)

    client := http.Client{}

    payload, err:= json.Marshal(message)
    if err != nil {
        panic(err)
    }
    datastring := string(payload[:])
    fmt.Printf("datastring: %s\n", datastring)

    request, err := http.NewRequest("POST", "https://android.googleapis.com/gcm/send", strings.NewReader("{\"registration_ids\":[\"" + reg + "\"], \"data\" : " + datastring + "}"))
    if err != nil {
        panic(err)
    }
    request.Header.Add("Content-Type","application/json")
    request.Header.Add("Authorization","key=AIzaSyBXM0IhHVTxedla2SR7h5v8aydliLNo6y4")
    resp, err2 := client.Do(request)
    if err2 != nil {
        panic(err2)
    }
    if resp != nil {
        fmt.Printf("responseCode: %d\n", resp.StatusCode);
    }
}

func (c App) Receive(receiver, num_from, messageReceived string) revel.Result {
    fmt.Printf("received text message")
    conduit := room.GetConduit(receiver)
    if conduit == nil {
        fmt.Printf("conduit was null when receiving message")
        panic("dang")
    }
    sms := room.SMSMessage{Num : num_from, Message : messageReceived}
    conduit.Received <- sms
    fmt.Printf("received, sent on channel\n")
    c.Response.Status = 201
    return c.RenderText("")
}

func (c App) Room(num string) revel.Result {
    return c.Render(num)
}

func (c App) Home() revel.Result {
    return c.RenderText("woooooooooo")
}

func (c App) Websock(num string, ws *websocket.Conn) revel.Result {
    // In order to select between websocket messages and subscription events, we
    // need to stuff websocket events into a channel.
    fmt.Printf("ws num: %s\n", num)
    conduit := room.GetConduit(num)
    if conduit == nil {
        fmt.Printf("ws conduit was null")
        panic("woah")
    }
    messagesToSend := make(chan room.SMSMessage)
    go func() {
        var msg string
        for {
            err := websocket.Message.Receive(ws, &msg)
            if err != nil {
                close(messagesToSend)
                return
            }
            fmt.Printf("got message from ws: %s\n", msg)
            var sms room.SMSMessage
            json.Unmarshal([]byte(msg), &sms)
            messagesToSend <- sms
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
