package websocket_routing

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"smswebproxy/app/appconnections"
	"smswebproxy/app/models"
	"strings"
)

func sendSmsMessage(message models.SMSMessage, appconn *appconnections.AppConnection, ws *websocket.Conn) bool {
	fmt.Printf("sendingMessage: %s\n", message.Message)

	client := http.Client{}

	payload, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	datastring := string(payload[:])
	fmt.Printf("datastring: %s\n", datastring)
	fmt.Printf("regid: %s\n", appconn.RegId)
	fmt.Printf("gcm: %s\n", os.Getenv("GCM_AUTH_KEY"))

	request, err := http.NewRequest("POST", "https://android.googleapis.com/gcm/send", strings.NewReader("{\"registration_ids\":[\""+appconn.RegId+"\"], \"data\" : "+datastring+"}"))
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
	return true
}
