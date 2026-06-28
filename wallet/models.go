package main

import (
	"fmt"
	"strconv"
	"time"

	u "gollux/utils"
)

type Account struct {
	ID          int       `json:"id"`
	AccountID   int       `json:"account_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	MobileNo    string    `json:"mobile_no"`
	Balance     float64   `json:"balance"`
	LastUpdated time.Time `json:"last_updated"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
	Remarks     string    `json:"remarks"`
	Number      int       `json:"number"`
}
type Transaction struct {
	ID              int       `json:"id"`
	AccountID       int       `json:"account_id"`
	RefNo           int       `json:"ref_no"`
	Description     string    `json:"description"`
	Amount          float64   `json:"amount"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Remarks         string    `json:"remarks"`
	CreatedAt       time.Time `json:"created_at"`
	PreviousBalance float64   `json:"previous_balance"`
	CurrentBalance  float64   `json:"current_balance"`
	SenderID        int       `json:"sender_id"`
}
type Balance struct {
	AccountID int     `json:"account_id"`
	Remarks   string  `json:"remarks"`
	Balance   float64 `json:"balance"`
	Status    bool    `json:"status"`
}

// Transaction types
const (
	TopUp      = "topup"
	Deposit    = "deposit"
	Withdrawal = "withdrawal"
	Transfer   = "transfer"
	Payment    = "payment"
	Refund     = "refund"
	Rollback   = "rollback"
)

func (trans *Transaction) Add() (*Transaction, bool) {
	acct := InquireAccount(trans.AccountID)
	if !acct.Status {
		trans.Status = "false"
		trans.Remarks = acct.Remarks
		return trans, false
	}
	if !trans.CanDoTransaction(acct.Balance) {
		trans.Status = "false"
		trans.Remarks = "Insufficient Balance"
		return trans, false
	}
	trans.PreviousBalance = acct.Balance
	trans.CreatedAt = time.Now()
	trans.Status = "SUCCESS"
	trans.Remarks = trans.Type
	errdb := DBM.Insert(trans)
	if errdb != nil {
		//panic(errdb)
		trans.Status = "false"
		trans.Remarks = "Error adding transaction!"
		return trans, false
	}
	trans.UpdateBalance()
	return trans, true
}
func (entry *Transaction) CanDoTransaction(balance float64) bool {
	if entry.Type == Payment || entry.Type == Transfer || entry.Type == Withdrawal {
		if balance >= entry.Amount {
			return true
		} else {
			return false
		}
	}
	return true
}
func (entry *Transaction) UpdateBalance() bool {
	var err error
	if entry.Type == Payment || entry.Type == Transfer || entry.Type == Withdrawal {
		_, err = DBM.Exec(`update accounts set last_updated=?, balance=coalesce(balance,0) - ?  where account_id=?`, time.Now(), entry.Amount, entry.AccountID)
	} else if entry.Type == TopUp || entry.Type == Deposit || entry.Type == Refund || entry.Type == Rollback {
		_, err = DBM.Exec(`update accounts set last_updated=?, balance=coalesce(balance,0) + ?  where account_id=?`, time.Now(), entry.Amount, entry.AccountID)
	}
	if err != nil {
		panic(err)
		return false
	}
	entry.SetCurrentBalance()
	return true
}

func (entry *Transaction) SetCurrentBalance() {
	var err error
	_, err = DBM.Exec(`update transactions set current_balance=account.balance from(select balance from accounts where account_id=?) as account  where transactions.id=?`, entry.AccountID, entry.ID)
	if err != nil {
		panic(err)
	}
}

func (acct *Account) NewAcount() (string, bool) {
	if resp, ok := acct.Validate(); !ok {
		return resp, ok
	}
	number, err := strconv.Atoi(u.GenNumCode(9))
	if err != nil {
		fmt.Println("Conversion error:", err)
	}
	acct.Balance = 0
	acct.CreatedAt = time.Now()
	acct.LastUpdated = time.Now()
	acct.Status = "Active"
	acct.Remarks = "Ok"
	acct.Number = number
	errdb := DBM.Insert(acct)
	if errdb != nil {
		panic(errdb)
		return "Error on new account creation", false
	}

	return "Account created!", true
}
func InquireAccount(AccountID int) *Balance {
	var p Account
	var balance Balance
	_, err := DBM.Query(&p, `SELECT * FROM accounts where account_id=?`, AccountID)
	if err != nil {
		panic(err)
		return &balance
	}
	if p.ID != 0 {
		balance.AccountID = p.ID
		balance.Status = true
		balance.Balance = p.Balance
		balance.Remarks = p.Remarks
		return &balance
	}
	balance.AccountID = p.ID
	balance.Status = false
	balance.Balance = 0
	balance.Remarks = "Account not exist!"
	return &balance
}

func (acct *Account) Validate() (string, bool) {
	var p Account
	_, err := DBM.Query(&p, `SELECT * FROM accounts where account_id=?`, acct.AccountID)
	if err != nil {
		panic(err)
		return "Account error validation", false
	}
	if p.ID != 0 {
		return "Account already exist!", false
	}
	return "Passed", true
}
