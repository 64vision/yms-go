package main

import (
	"encoding/json"
	"fmt"
	u "gollux/utils"
	"net/http"
	"time"
)

type Webhook struct {
	Event     string      `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
	Data      PaymentData `json:"data"`
}

type PaymentData struct {
	PaymentID         string    `json:"payment_id"`
	Status            string    `json:"status"`
	Amount            int       `json:"amount"`
	Currency          string    `json:"currency"`
	Gateway           string    `json:"gateway"`
	Reference         string    `json:"reference"`
	ProviderReference string    `json:"provider_reference"`
	PaidAt            time.Time `json:"paid_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func ApiWebhook(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	entry := &Webhook{}

	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(entry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	fmt.Print(entry.Data.Reference)

	u.Respond(w, resp)

}
