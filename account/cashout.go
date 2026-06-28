package account

import (
	"encoding/json"
	"net/http"
	"time"

	"gollux/sms"
	u "gollux/utils"
)

type Cashout struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	Fee       float64   `json:"fee"`
	ReqAt     time.Time `json:"req_at"`
	ProcessAt time.Time `json:"process_at"`
	ProcessBy int       `json:"process_by"`
	Bank      Bank      `json:"bank"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Remarks   string    `json:"remarks"`
}
type CashoutList struct {
	ID          int       `json:"id"`
	AccountID   int       `json:"account_id"`
	AccountName string    `json:"account_name"`
	Fee         float64   `json:"fee"`
	ReqAt       time.Time `json:"req_at"`
	ProcessAt   time.Time `json:"process_at"`
	ProcessBy   int       `json:"process_by"`
	Bank        Bank      `json:"bank"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Remarks     string    `json:"remarks"`
}
type Bank struct {
	Name        string `json:"name"`
	AccountName string `json:"account_name"`
	Number      string `json:"number"`
	Type        string `json:"type"`
}

var GetCashout = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	cashout := &Cashout{}
	err := json.NewDecoder(r.Body).Decode(cashout) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	cashout.AccountID = r.Context().Value("user").(int)
	resp := cashout.Get()
	u.Respond(w, resp)
}

var UpdateCashout = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	cashout := &Cashout{}
	err := json.NewDecoder(r.Body).Decode(cashout) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	cashout.ProcessBy = r.Context().Value("user").(int)
	resp := cashout.Update()
	u.Respond(w, resp)
}

var DoCashout = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	cashout := &Cashout{}
	err := json.NewDecoder(r.Body).Decode(cashout) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	cashout.AccountID = r.Context().Value("user").(int)
	resp := cashout.Request()
	u.Respond(w, resp)
}

func (cashout *Cashout) Request() map[string]interface{} {
	if cashout.Amount < 100 {
		return u.Message(false, "Minimum 100 for cashout!.")
	}
	if !cashout.CanCashout() {
		return u.Message(false, "Still processing your pending cashout.")
	}
	acct := InquireAccount(cashout.AccountID)

	if acct.Balance < cashout.Amount {
		return u.Message(false, "Insufficient credits to cash!")
	}
	cashout.Status = "Requested"
	cashout.ReqAt = time.Now()
	errdb := DBM.Insert(cashout)
	if errdb != nil {
		panic(errdb)
	}
	trans := Transaction{}

	trans.AccountID = cashout.AccountID
	trans.Amount = cashout.Amount
	trans.PreviousBalance = acct.Balance
	trans.CreatedAt = time.Now()
	trans.Type = Withdrawal
	trans.RefNo = cashout.ID
	trans.Status = "SUCCESS"
	trans.Description = "Cashout Request"
	trans.Remarks = Withdrawal
	errdb = DBM.Insert(&trans)
	if errdb != nil {
		panic(errdb)
	}
	trans.UpdateBalance()
	sms.Send("09177723286", "Hi admin, New cashout has been requested.")
	sms.Send("09156033392", "Hi admin, New cashout has been requested.")
	response := u.Message(true, "Cashout request submitted successfully!")
	return response
}

func (qry *Query) CashoutQry() map[string]interface{} {
	var items []CashoutList
	_, errdb := DBM.Query(&items, qry.Query)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, "Ok")
	response["items"] = items
	return response
}
func (cashout *Cashout) CanCashout() bool {
	var item Cashout
	res, errdb := DBM.Query(&item, `select * from cashouts where status=? and account_id=? `, "Requested", cashout.AccountID)
	if errdb != nil {
		panic(errdb)
	}
	if res.RowsReturned() > 0 {
		return false
	}

	return true
}

func (cashout *Cashout) Get() map[string]interface{} {
	var item Cashout
	res, errdb := DBM.Query(&item, `select * from cashouts where account_id=? and status=?`, cashout.AccountID, "Requested")
	if errdb != nil {
		panic(errdb)
	}
	if res.RowsReturned() == 0 {
		return u.Message(false, "No Data Available")
	}
	response := u.Message(true, "Ok")
	response["item"] = item
	return response
}
func (cashout *Cashout) Update() map[string]interface{} {

	var item Cashout
	res, errdb := DBM.Query(&item, `select * from cashouts where id=?`, cashout.ID)
	if errdb != nil {
		panic(errdb)
	}
	if res.RowsReturned() == 0 {
		return u.Message(false, "Cashout not found.")
	}
	if item.Status == "Cancelled" {
		return u.Message(false, "This request was already cancelled!")
	}
	item.Remarks = cashout.Remarks
	item.Status = cashout.Status
	errdb = DBM.Update(&item)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	if cashout.Status == "Cancelled" {
		trans := Transaction{}
		acct := InquireAccount(item.AccountID)

		trans.AccountID = item.AccountID
		trans.Amount = item.Amount
		trans.PreviousBalance = acct.Balance
		trans.CreatedAt = time.Now()
		trans.Type = Deposit
		trans.RefNo = item.ID
		trans.Status = "SUCCESS"
		trans.Description = "Cashout Cancelled"
		trans.Remarks = Deposit
		errdb = DBM.Insert(&trans)
		if errdb != nil {
			panic(errdb)
		}
		trans.UpdateBalance()
	}
	if cashout.Status == "Completed" {
		acct, _ := GetAccountByID(cashout.AccountID)
		sms.Send(acct.Username, "Your cashout request has been successful. Enjoy BLAZER!")
	}

	response := u.Message(true, cashout.Status)
	return response
}
