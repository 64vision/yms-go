package main

import (
	"encoding/json"
	"fmt"
	"gollux/account"
	u "gollux/utils"
	"net/http"
	"time"
	"zerasuite/bookings"
)

func WithdrawContainer(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	entry := &bookings.Booking{}
	bal := &account.Balance{}
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(entry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	entry.ClientID = r.Context().Value("user").(int)
	bal = account.InquireAccount(entry.ClientID)
	if bal.Balance < entry.WithdrawalAmount {
		u.Respond(w, u.Message(false, "Insufficient funds! Please Top UP."))
		return
	}
	resp = entry.Withdraw()
	if resp["status"].(bool) {
		fmt.Println("W Status", resp["status"].(bool))
		trans := &account.Transaction{}
		trans.Type = "payment"
		trans.AccountID = entry.ClientID
		trans.Amount = entry.WithdrawalAmount
		trans.CreatedAt = time.Now()
		trans.RefNo = entry.ID
		trans.Description = "Container withdrawal"
		trans.Add()

	}
	u.Respond(w, resp)
}

func NewBooking(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	entry := &bookings.Booking{}
	bal := &account.Balance{}
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(entry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	entry.ClientID = r.Context().Value("user").(int)
	bal = account.InquireAccount(entry.ClientID)
	if bal.Balance < entry.TotalAmount {
		u.Respond(w, u.Message(false, "Insufficient funds! Please Top UP."))
		return
	}
	resp = entry.AddBooking()
	if resp["status"].(bool) {
		fmt.Println("Booking Status", resp["status"].(bool))
		trans := &account.Transaction{}
		trans.Type = "payment"
		trans.AccountID = entry.ClientID
		trans.Amount = entry.TotalAmount
		trans.CreatedAt = time.Now()
		trans.RefNo = resp["booking_id"].(int)
		trans.Description = "Yard slot booking"
		trans.Add()

	}
	u.Respond(w, resp)
}
func ViewBooking(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	entry := &bookings.Booking{}
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(entry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp = entry.GetBookingByID()
	u.Respond(w, resp)
}
