package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"hyperball.com/account"
	"hyperball.com/sms"
	u "hyperball.com/utils"
)

type Message struct {
	ID           int       `json:"id"`
	Receiver     string    `json:"receiver"`
	ReceiverType string    `json:"receiver_type"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
	MessageTxt   string    `json:"message_txt"`
	CreatedBy    int       `json:"created_by"`
	Status       string    `json:"status"`
}

var ListMessages = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	resp := MessageList()
	u.Respond(w, resp)
}
var SendMessage = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	data := &Message{}
	err := json.NewDecoder(r.Body).Decode(data) //decode the request body into struct and failed if any error occu
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	data.CreatedBy = r.Context().Value("user").(int)
	fmt.Println(data)
	resp := data.Send()
	u.Respond(w, resp)
}

func (msg *Message) Send() map[string]interface{} {
	msg.CreatedAt = time.Now()
	msg.Status = "pending"
	errdb := DBM.Insert(msg)
	if errdb != nil {
		fmt.Println(errdb.Error())
		panic(errdb.Error())
		return u.Message(false, "Failed")
	}
	if msg.ReceiverType == "Player" {
		sms.Send(msg.Receiver, msg.MessageTxt)
	} else if msg.ReceiverType == "All Players" {
		sendToAll(msg.MessageTxt)
	}

	return u.Message(true, "Success!")

}
func sendToAll(msg string) {
	var accts []account.Account
	_, err := account.DBM.Query(&accts, `select * from accounts`)
	if err != nil {
		panic(err)
	}
	for _, player := range accts {
		sms.Send(player.MobileNo, msg)
	}

}
func MessageList() map[string]interface{} {
	var msgs []Message
	_, err := DBM.Query(&msgs, `select * from messages order by id desc`)
	if err != nil {
		panic(err)
	}

	resp := u.Message(true, "Ok")
	resp["items"] = msgs
	return resp

}
