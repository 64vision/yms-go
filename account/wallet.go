package account

import (
	"fmt"
	"time"

	"gollux/sms"
	u "gollux/utils"
)

type Transaction struct {
	ID              int                    `json:"id"`
	AccountID       int                    `json:"account_id"`
	RefNo           int                    `json:"ref_no"`
	Description     string                 `json:"description"`
	Amount          float64                `json:"amount"`
	Type            string                 `json:"type"`
	Status          string                 `json:"status"`
	Remarks         string                 `json:"remarks"`
	CreatedAt       time.Time              `json:"created_at"`
	PreviousBalance float64                `json:"previous_balance"`
	CurrentBalance  float64                `json:"current_balance"`
	SenderID        int                    `json:"sender_id"`
	Details         map[string]interface{} `json:"details"`
}
type Balance struct {
	AccountID int     `json:"account_id"`
	Remarks   string  `json:"remarks"`
	Balance   float64 `json:"balance"`
	Status    bool    `json:"status"`
}

type Request struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	AccountID int       `json:"account_id"`
	ReqAt     time.Time `json:"req_at"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Code      string    `json:"code"`
	Receiver  string    `json:"receiver"`
	Message   string    `json:"message"`
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

func (req *Request) NewTopUp() map[string]interface{} {

	receiver, _ := GetAccountByMobile(req.Receiver)
	if receiver.ID == 0 {
		return u.Message(false, "Mobile no. is not yet registered.")
	}
	req.ReqAt = time.Now()
	req.Status = "OK"
	req.Code = u.GenNumCode(4)

	errdb := DBM.Insert(req)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	} else {

		var receiverTrans Transaction
		receiverTrans.Amount = req.Amount
		receiverTrans.AccountID = receiver.ID
		receiverTrans.Type = req.Type
		receiverTrans.Description = "Received Credits from  Blazing Sphere"
		receiverTrans.Remarks = req.Type
		receiverTrans.Create()
		sms.Send(req.Receiver, req.Message)
	}
	resp := u.Message(true, "Successful!")
	resp["requested"] = req
	return resp
}

func (req *Request) Add() map[string]interface{} {
	sender, _ := GetAccountByID(req.AccountID)
	if sender.Level < 2 {
		return u.Message(false, "No permission to do this transaction")
	}
	if sender.Balance < req.Amount {
		return u.Message(false, "Insufficient credits to transfer!")
	}
	receiver, _ := GetAccountByMobile(req.Receiver)
	if receiver.ID == 0 {
		return u.Message(false, "Mobile no. is not yet registered.")
	}
	req.ReqAt = time.Now()
	req.Status = "For Verification"
	req.Code = u.GenNumCode(4)
	errdb := DBM.Insert(req)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	} else {

		_msg := "You requested " + fmt.Sprintf("%.2f", req.Amount) + " to mobile no. " + req.Receiver + ". To proceed here is your OTP " + req.Code + "."
		sms.Send(sender.MobileNo, _msg)
	}
	resp := u.Message(true, "Request Successful!")
	resp["requested"] = req
	return resp
}

func (req *Request) Proceed() map[string]interface{} {
	var p Request
	_, err := DBM.Query(&p, `SELECT * FROM requests where id=?`, req.ID)
	if err != nil {
		return u.Message(false, err.Error())
	}
	if p.Status != "For Verification" {
		return u.Message(false, "Request was already completed!")
	}
	if p.Code != req.Code {
		return u.Message(false, "Incorrect OTP!")
	}
	p.Status = "Completed"
	errdb := DBM.Update(&p)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	} else {
		sender, _ := GetAccountByID(p.AccountID)

		var senderTrans Transaction
		senderTrans.AccountID = p.AccountID
		senderTrans.Amount = p.Amount
		senderTrans.Type = Transfer
		senderTrans.Description = "Send Credits to " + p.Receiver
		senderTrans.Remarks = Transfer
		senderTrans.Create()

		receiver, _ := GetAccountByMobile(p.Receiver)
		var receiverTrans Transaction
		receiverTrans.Amount = p.Amount
		receiverTrans.AccountID = receiver.ID
		receiverTrans.Type = Deposit
		receiverTrans.Description = "Received Credits from " + sender.MobileNo
		receiverTrans.Remarks = Deposit
		receiverTrans.Create()

		_msg := "You have recieved " + fmt.Sprintf("%.2f", p.Amount) + " HBR coins from  " + sender.MobileNo + ". Let's Roll and Win!"
		fmt.Println("Receiver", sender.MobileNo)
		sms.Send(p.Receiver, _msg)

	}
	resp := u.Message(true, "Sending credits complete!")
	return resp
}

func (trans *Transaction) Create() map[string]interface{} {
	acct := InquireAccount(trans.AccountID)

	trans.PreviousBalance = acct.Balance
	trans.CreatedAt = time.Now()
	trans.Status = "SUCCESS"
	trans.Remarks = trans.Type
	errdb := DBM.Insert(trans)
	if errdb != nil {
		panic(errdb)
		return u.Message(false, errdb.Error())
	}
	trans.UpdateBalance()
	resp := u.Message(true, "Successfully!")
	resp["transaction"] = trans
	return resp
}

func (trans *Transaction) Add() map[string]interface{} {
	acct := InquireAccount(trans.AccountID)
	fmt.Println(acct)
	if !acct.Status {
		return u.Message(false, "Can't add transaction!")
	}
	if !trans.CanDoTransaction(acct.Balance) {
		return u.Message(false, "Insufficient Balance!")
	}
	trans.PreviousBalance = acct.Balance
	trans.CreatedAt = time.Now()
	trans.Status = "SUCCESS"
	trans.Remarks = trans.Type
	errdb := DBM.Insert(trans)
	if errdb != nil {
		panic(errdb)
		return u.Message(false, errdb.Error())
	}
	trans.UpdateBalance()
	resp := u.Message(true, "Successfully!")
	resp["transaction"] = trans
	return resp
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
		_, err = DBM.Exec(`update accounts set balance=coalesce(balance,0) - ?, balance_update_at=?  where id=?`, entry.Amount, time.Now(), entry.AccountID)
	} else if entry.Type == TopUp || entry.Type == Deposit || entry.Type == Refund || entry.Type == Rollback {
		_, err = DBM.Exec(`update accounts set  balance=coalesce(balance,0) + ?, balance_update_at=?  where id=?`, entry.Amount, time.Now(), entry.AccountID)
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
	_, err = DBM.Exec(`update transactions set current_balance=account.balance from(select balance from accounts where id=?) as account  where transactions.id=?`, entry.AccountID, entry.ID)
	if err != nil {
		panic(err)
	}
}

func InquireAccount(AccountID int) *Balance {
	var p Account
	var balance Balance
	_, err := DBM.Query(&p, `SELECT * FROM accounts where id=?`, AccountID)
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

func (acct *Account) GetTransactions() map[string]interface{} {
	var transactions []Transaction
	_, err := DBM.Query(&transactions, `SELECT * FROM transactions where account_id=? order by id desc limit 50`, acct.ID)
	if err != nil {
		return u.Message(false, err.Error())
	}
	resp := u.Message(true, "Result!")
	resp["transactions"] = transactions
	return resp

}
