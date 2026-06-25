package account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hyperball.com/sms"
	u "hyperball.com/utils"
)

type CallbackRes struct {
	AcquirementId     string      `json:"acquirementId"`
	MerchantTransId   string      `json:"merchantTransId"`
	AcquirementStatus string      `json:"acquirementStatus"`
	TransactionId     string      `json:"transactionId"`
	FinishedTime      string      `json:"finishedTime"`
	CreatedTime       string      `json:"createdTime"`
	OrderAmount       interface{} `json:"orderAmount"`
}
type UrlRes struct {
	Url string `json:"url"`
}

type OrderGcashRes struct {
	Status          bool       `json:"status"`
	Code            int        `json:"code"`
	Message         string     `json:"message"`
	AcquirementId   string     `json:"acquirementId"`
	CheckoutUrl     string     `json:"checkoutUrl"`
	MerchantTransId string     `json:"merchantTransId"`
	ResultInfo      ResultInfo `json:"resultInfo"`
	TransactionId   string     `json:"transactionId"`
}

type ResultInfo struct {
	ResultCodeId string `json:"resultCodeId"`
	ResultMsg    string `json:"resultMsg"`
	ResultStatus string `json:"resultStatus"`
	ResultCode   string `json:"resultCode"`
}
type BuyCreditList struct {
	ID          int        `json:"id"`
	TransID     string     `json:"trans_id"`
	UserID      int        `json:"user_id"`
	AccountName string     `json:"account_name"`
	Request     CreditsReq `json:"request"`
	ResponseUrl string     `json:"response_url"`
	Status      string     `json:"status"`
	ReqDt       time.Time  `json:"req_dt"`
	Partner     string     `json:"partner"`
	Amount      float64    `json:"amount"`
	ReqAmount   float64    `json:"req_amount"`
}
type BuyCredit struct {
	ID          int        `json:"id"`
	TransID     string     `json:"trans_id"`
	UserID      int        `json:"user_id"`
	Request     CreditsReq `json:"request"`
	ResponseUrl string     `json:"response_url"`
	Status      string     `json:"status"`
	ReqDt       time.Time  `json:"req_dt"`
	Partner     string     `json:"partner"`
	Amount      float64    `json:"amount"`
	ReqAmount   float64    `json:"req_amount"`
}
type CreditsReq struct {
	Name            string      `json:"name"`
	Amount          float64     `json:"amount"`
	TransactionCode string      `json:"transaction_code"`
	Transaction     string      `json:"transaction"`
	Email           string      `json:"email"`
	MobileNo        string      `json:"mobile_no"`
	Callback        string      `json:"callback"`
	Collection      CreditsColl `json:"collection"`
}
type CreditsColl struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Code   string  `json:"code"`
}
type MayaCreateRes struct {
	CheckoutId      string `json:"checkout_id"`
	TransactionCode string `json:"transaction_code"`
	CheckoutUrl     string `json:"checkout_url"`
}

type MayaRes struct {
	Status          string    `json:"status"`
	CheckoutId      string    `json:"checkout_id"`
	TransactionCode string    `json:"transaction_code"`
	MerchantRefno   string    `json:"merchant_refno"`
	Amount          string    `json:"amount"`
	Timestamp       time.Time `json:"timestamp"`
}

// GCASH
const (
	URL                = "https://payment-master-v1.surepay-prod.com/create_order"
	CLIENTID           = "sp_pwp80194cd74006b693738131c4ecd91"
	CLIENTKEY          = "sp_pwp80194cd74006b693738131c4ecd91"
	PARTNER            = "GCASH"
	COMPLETEURL        = "https://web.servehpbr.com/payment.html"
	CANCELURL          = "https://web.servehpbr.com/failed.html"
	GCASH_MERCHANTCODE = "HYPERBALL"
	GCASH_MERCHANTNAME = "HYPERBALL"
	MAYAURL            = "https://maya-prod.surepayinc.com.ph/api/checkout"
	MAYAURLCALBACK     = "https://acct.servehpbr.com/credits/maya_callback"
	MAYACLIENTSECRETE  = "JQPFSDH4YfO7bQfG3Py8DN3jGO48acn0Gaz9cnQtbT4SMGGuXnDyEbtgzPjmrwWx"
	MAYACLIENTID       = "872cd1f8-b8cf-4eff-a829-8055c727e93a"
)

var DoBuyCredits = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	data := &BuyCredit{}
	err := json.NewDecoder(r.Body).Decode(data) //decode the request body into struct and failed if any error occu
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	data.UserID = r.Context().Value("user").(int)
	fmt.Println(data)
	resp := data.Buy()
	u.Respond(w, resp)
}
var Callback = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	model := &CallbackRes{}
	err := json.NewDecoder(r.Body).Decode(model) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	model.BuyComplete()
	resp := u.Message(true, "OK!")
	u.Respond(w, resp)

}
var MayaCallback = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	model := &MayaRes{}
	err := json.NewDecoder(r.Body).Decode(model) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	model.Complete()
	resp := u.Message(true, "OK!")
	u.Respond(w, resp)

}

func (qry *Query) BuyCreditsQry() map[string]interface{} {
	fmt.Println("BuyCreditsQry")
	var items []BuyCreditList
	_, errdb := DBM.Query(&items, qry.Query)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, "Ok")
	response["items"] = items
	return response
}
func (entry *CallbackRes) BuyComplete() {
	fmt.Println("CallbackRes Complete")

	res, err := DBM.Exec(`UPDATE buy_credits set status=? where trans_id=?`, entry.AcquirementStatus, entry.MerchantTransId)
	if err != nil {
		panic(err)
	}
	if entry.AcquirementStatus == "SUCCESS" {
		processCredits(entry.MerchantTransId)
	}
	fmt.Println("CallbackRes Complete", res)

}

func (entry *MayaRes) Complete() {
	fmt.Println("Maya CallbackRes Complete")
	if entry.Status == "S" {
		entry.Status = "SUCCESS"
	}

	res, err := DBM.Exec(`UPDATE buy_credits set status=? where trans_id=?`, entry.Status, entry.MerchantRefno)
	if err != nil {
		panic(err)
	}
	if res.RowsAffected() > 0 {
		if entry.Status == "SUCCESS" {
			processCredits(entry.MerchantRefno)
		}

	}

}

func processCredits(transaction_code string) {
	var req BuyCredit
	trans := Transaction{}
	_, err := DBM.Query(&req, `select * from buy_credits where trans_id=?`, transaction_code)
	if err != nil {
		panic(err)
	}

	acct, _ := GetAccountByID(req.UserID)

	trans.AccountID = req.UserID
	trans.Amount = req.Amount
	trans.PreviousBalance = acct.Balance
	trans.CreatedAt = time.Now()
	trans.Type = TopUp
	trans.RefNo = req.ID
	trans.Status = "SUCCESS"
	trans.Description = "Buy credits via " + req.Partner
	trans.Remarks = req.Partner
	errdb := DBM.Insert(&trans)
	if errdb != nil {
		panic(errdb)
	}
	trans.UpdateBalance()

	_msg := fmt.Sprintf("%.2f", req.Amount) + " HBR credits was successfully added to your account. Let's roll and win!"
	sms.Send(acct.MobileNo, _msg)

}

func (buy *BuyCredit) Buy() map[string]interface{} {
	merchantTransId := "shbr" + strconv.Itoa(buy.UserID) + "-" + genTransCode()
	//account := GetUserMeta(buy.UserID)
	buy.Amount = buy.ReqAmount
	buy.TransID = merchantTransId
	//buy.Request = *reqData
	buy.Status = "Pending"
	buy.ReqDt = time.Now()
	errdb := DBM.Insert(buy)
	if errdb != nil {
		fmt.Println(errdb.Error())
		panic(errdb.Error())
		return u.Message(false, "Failed")
	}
	var payRes map[string]interface{}
	if buy.Partner == "MAYA" {
		_amount := fmt.Sprintf("%f", buy.ReqAmount)
		fmt.Println("maya amount", _amount)
		acct, _ := GetAccountByID(buy.UserID)
		payRes = MayaCreateOrder(merchantTransId, _amount, acct.FirstName, acct.LastName, "64marbles@gmail.com")
	} else if buy.Partner == "GCASH" {
		payRes = GcashCreateOrder(merchantTransId, buy.ReqAmount)
	}
	return payRes
}

func MayaCreateOrder(merchantTransId string, amount string, firstname string, lastname string, email string) map[string]interface{} {

	var jsonData = `{
		"firstname": "` + firstname + `",
		"lastname": "` + lastname + `",
		"email": "` + email + `",
		"amount": "` + amount + `",
		"callback_success": "` + COMPLETEURL + `",
		"callback_cancelled":  "` + CANCELURL + `",
		"callback_failed":  "` + CANCELURL + `",
		"merchant_refno": "` + merchantTransId + `",
		"callback_url": "` + MAYAURLCALBACK + `"
	}`
	//fmt.Println("jsonData:", jsonData)
	jsonStr := strings.NewReader(jsonData)
	//fmt.Println("jsonStr:", jsonStr)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", MAYAURL, jsonStr)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Sp-Client-Id", MAYACLIENTID)
	req.Header.Set("Sp-Client-Secret", MAYACLIENTSECRETE)
	// Make HTTP POST request and return message SID
	response := u.Message(true, "Successful")
	resp, errb := client.Do(req)

	if errb != nil {
		fmt.Println("Error:", errb)

	}
	defer resp.Body.Close()

	// Read the response body and convert it to a string
	_body, errc := ioutil.ReadAll(resp.Body)
	if errc != nil {
		fmt.Println("Error reading response body:", errc)

	}

	// Convert the response body to a string
	responseString := string(_body)

	fmt.Println("HTTP Status Code:", resp.Status)
	fmt.Println("Response Body as String:", responseString)
	fmt.Println("------")
	var data MayaCreateRes
	fmt.Println(resp.Body)
	// Decode the response body into a struct

	errr := json.Unmarshal(_body, &data)
	if errr != nil {
		fmt.Println("Error decoding response body:", errr)

	}

	// Now you can work with the decoded struct
	if resp.Status >= "201 Created" {
		fmt.Println("pass here2")
		response = u.Message(true, "OK")
		response["response"] = data
		return response
	}
	return response
}

func GcashCreateOrder(merchantTransId string, amount float64) map[string]interface{} {
	str := strconv.FormatFloat(amount, 'f', -1, 64)
	total := str + "00"
	var jsonData = `{
		"req_name": "order",
		"merchantCode": "` + GCASH_MERCHANTCODE + `",
		"merchantName": "` + GCASH_MERCHANTNAME + `",
		"order": { 
				"merchantTransId": "` + merchantTransId + `", 
				"merchantTransType": "CASHIER_ORDER", 
				"orderMemo": "Surepay FEES", 
				"orderTitle": "Surepay of LGU", 
				"seller":{ 
				  "userId":"Surepay", 
				  "externalUserId":"Surepay", 
				  "externalUserType":"Surepay" 
				  }, 
				  "buyer":{ 
					"userId":"", 
					"externalUserId":"Surepay", 
					"externalUserType":"Surepay" 
				  }, 
				  
				"orderAmount": { "currency": "PHP", "value": "` + total + `" }, 
				"createdTime": "", 
				"expiryTime": "" 
			},
		"notificationUrls": [ 
				{ "type": "PAY_RETURN", "url": "` + COMPLETEURL + `" }, 
				{ "type": "CANCEL_RETURN", "url": "` + CANCELURL + `" }
			]   
	}`
	//{ "type": "PAY_RETURN", "url": "https://guest.eboracay.net/complete" },
	//{ "type": "CANCEL_RETURN", "url": "https://guest.eboracay.net/cancelled" }
	jsonStr := strings.NewReader(jsonData)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", URL, jsonStr)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("surepaykey", CLIENTKEY)
	// Make HTTP POST request and return message SID
	response := u.Message(true, "Successful")
	resp, errb := client.Do(req)

	if errb != nil {
		fmt.Println("Error:", errb)

	}
	defer resp.Body.Close()

	// Read the response body and convert it to a string
	_body, errc := ioutil.ReadAll(resp.Body)
	if errc != nil {
		fmt.Println("Error reading response body:", errc)

	}

	// Convert the response body to a string
	responseString := string(_body)

	fmt.Println("HTTP Status Code:", resp.Status)
	fmt.Println("Response Body as String:", responseString)
	fmt.Println("------")
	var data OrderGcashRes
	fmt.Println(resp.Body)
	// Decode the response body into a struct

	errr := json.Unmarshal(_body, &data)
	if errr != nil {
		fmt.Println("Error decoding response body:", errr)

	}

	// Now you can work with the decoded struct
	fmt.Println(data.Status)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		//fmt.Println(resp.Body)
		// err := json.NewDecoder(resp.Body).Decode(&data)
		// if err != nil {
		// 	panic(err.Error())
		// 	return u.Message(false, "Failed here!")
		// }
		if !data.Status {
			fmt.Println("pass pause")
			return u.Message(false, data.Message)
		}
		fmt.Println("pass here2")
		response = u.Message(true, "OK")
		response["response"] = data
		return response
	}
	return response
}
func genTransCode() string {
	// min := 1000000
	// max := 9000000
	// return strconv.Itoa(min + rand.Intn(max-min))
	_now := time.Now()
	formattedTime := _now.Format("060102150405")
	return formattedTime
}
