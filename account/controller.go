package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"hyperball.com/sms"
	u "hyperball.com/utils"
)

func AddRequest(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Request{}
	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	qry.AccountID = r.Context().Value("user").(int)
	if qry.Type == "topup" {
		resp := qry.NewTopUp()
		u.Respond(w, resp)
	} else {
		resp := qry.Add()
		u.Respond(w, resp)
	}

}
func ProceedRequest(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Request{}
	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	qry.AccountID = r.Context().Value("user").(int)
	resp := qry.Proceed()
	u.Respond(w, resp)
}

func PlayerStats(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	resp := AccountStats()
	u.Respond(w, resp)
}

func BalanceInquire(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Account{}

	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	bal := InquireAccount(qry.ID)
	if bal.AccountID == 0 {
		u.Respond(w, u.Message(false, "Account not found!"))
		return
	}
	resp := u.Message(true, "Sent!")
	resp["balance"] = bal
	u.Respond(w, resp)
}

func PlayerCustomQry(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Query{}
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	if r.Context().Value("user").(int) != qry.AccountID {
		u.Respond(w, u.Message(false, "Permission denied!"))
		return
	}

	if qry.Type == "query" {
		resp = qry.CustomQry()
	} else if qry.Type == "update account" {
		prohibiteds := []string{"balance", "delete", "level", "truncate", "commission", "agent", "created_at", "remarks"}
		for _, prohibited := range prohibiteds {
			if strings.Contains(strings.ToLower(qry.Query), strings.ToLower(prohibited)) {
				u.Respond(w, u.Message(false, "Fuck you!"))
				return
			}
		}
		resp = qry.AccountUpdate()
	}

	u.Respond(w, resp)
}
func CustomQry(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Query{}
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	fmt.Println(qry.Type)
	if qry.Type == "query" {
		resp = qry.CustomQry()
	} else if qry.Type == "update account" {
		resp = qry.AccountUpdate()
	} else if qry.Type == "buycredits" {
		resp = qry.BuyCreditsQry()
	} else if qry.Type == "cashouts" {
		resp = qry.CashoutQry()
	}

	u.Respond(w, resp)
}

func Test(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Register modoule!", name)
	return message
}
func GetPlayersLocation(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	resp := PlayersLocation()
	u.Respond(w, resp)
}
func TestAccess(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	resp := u.Message(true, "OK")
	resp["account"] = r.Context().Value("user").(int)
	u.Respond(w, resp)
}
func LoginByPassword(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	account := &LogRequest{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	fmt.Println(account)
	resp := account.Login()
	u.Respond(w, resp)
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	admin := &Administrator{}
	err := json.NewDecoder(r.Body).Decode(admin) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := admin.Login()
	u.Respond(w, resp)
}

func Verify(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	verify := &VerifyReq{}
	err := json.NewDecoder(r.Body).Decode(verify) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := verify.Verify()
	u.Respond(w, resp)
}

func Register(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	account := &Registration{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	fmt.Println(account)
	resp := account.Add()
	u.Respond(w, resp)
}

func ForgotCode(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	acct := &Account{}
	err := json.NewDecoder(r.Body).Decode(acct) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if !u.IsValidPHMobileNumber(acct.MobileNo) {
		u.Message(false, "Invalid Mobile number!")
		return
	}
	account, _ := GetAccountByMobile(acct.MobileNo)
	if account.ID == 0 {
		u.Message(false, "Mobile number not found!!")
		return
	}
	code := u.GenNumCode(4)
	_hash := u.Md5hash(code)
	var qry Query
	qry.Table = "accounts"
	qry.Query = "update accounts set password='" + _hash + "' where mobile_no='" + acct.MobileNo + "'"
	qry.AccountUpdate()
	smsmsg := "Your temporary password: " + code + ". Please change your password once you logged in."
	send, resmsg := sms.Send(acct.MobileNo, smsmsg)
	if !send {
		u.Respond(w, u.Message(false, resmsg))
		return
	}
	resp := u.Message(true, "Sent!")
	u.Respond(w, resp)
}
func ResendCode(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	acct := &Account{}
	err := json.NewDecoder(r.Body).Decode(acct) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	account, _ := GetAccountByMobile(acct.MobileNo)
	if account.ID == 0 {
		u.Message(false, "Mobile number not found!!")
		return
	}
	code := u.GenNumCode(4)
	var qry Query
	qry.Table = "accounts"
	qry.Query = "update accounts set code='" + code + "' where mobile_no='" + acct.MobileNo + "'"
	qry.AccountUpdate()
	smsmsg := "Your verification code is " + code + ". Enter this code to verify your account."
	send, resmsg := sms.Send(acct.MobileNo, smsmsg)
	if !send {
		u.Respond(w, u.Message(false, resmsg))
		return
	}
	resp := u.Message(true, "Sent!")
	u.Respond(w, resp)
}
func GetTransactionsHistory(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	account := &Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.GetTransactions()
	u.Respond(w, resp)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	account := &Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if account.ID == 0 {
		u.Respond(w, u.Message(false, "Account ID is required"))
		return
	}
	acct, errAct := GetAccountByID(account.ID)
	if errAct != nil {
		u.Respond(w, u.Message(false, "Account not found"))
	}

	resp := u.Message(true, "Ok")
	resp["account"] = acct
	u.Respond(w, resp)
}
