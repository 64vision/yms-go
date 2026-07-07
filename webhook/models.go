package main

import (
	"encoding/json"
	"fmt"
	"gollux/account"
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
		panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	fmt.Println(entry.Data.Reference)
	if entry.Event == "payment.updated" {
		UpdateRequest(entry.Data.Reference, entry.Data.Status)
	}
	resp = u.Message(true, "OK")
	u.Respond(w, resp)

}
func UpdateRequest(reference string, status string) {
	res, err := DBM.Exec(`UPDATE buy_credits set status=? where trans_id=?`, status, reference)
	if err != nil {
		panic(err)
	}
	if res.RowsAffected() == 0 {
		fmt.Println("Transaction not found!")
		return
	}
	if status == "paid" {
		ProcessCredits(reference)
	}
	fmt.Println("Transaction OK")
}

func ProcessCredits(transaction_code string) {
	fmt.Println("ProcessCredits", transaction_code)
	var req account.BuyCredit
	trans := account.Transaction{}
	_, err := DBM.Query(&req, `select * from buy_credits where trans_id=?`, transaction_code)
	if err != nil {
		panic(err)
	}
	// fmt.Println("current Status", req.Status)
	// if req.Status == "paid" {
	// 	fmt.Println("Already Proccess")
	// 	return
	// }
	acct, _ := account.GetAccountByID(req.UserID)

	trans.AccountID = req.UserID
	trans.Amount = req.Amount
	trans.PreviousBalance = acct.Balance
	trans.CreatedAt = time.Now()
	trans.Type = "topup"
	trans.RefNo = req.ID
	trans.Status = "SUCCESS"
	trans.Description = "Buy credits via " + req.Partner + " - " + req.Gateway
	trans.Remarks = req.Partner
	errdb := DBM.Insert(&trans)
	if errdb != nil {
		panic(errdb)
	}
	trans.UpdateBalance()

}
