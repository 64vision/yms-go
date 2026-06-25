package main

import (
	"encoding/json"
	"net/http"

	"hyperball.com/account"
	u "hyperball.com/utils"
)

type QryParam struct {
	Type      string `json:"type"`
	Date      string `json:"date"`
	AccountID int    `json:"account_id"`
}

func Query(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &account.Query{}
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if qry.Type == "sales" {
		resp = qry.AgentSales()
	} else if qry.Type == "location" {
		resp = qry.AgentLocation()
	} else if qry.Type == "settlements" {
		resp = qry.Settlements()
	}
	u.Respond(w, resp)
}
func GetPerformance(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	qry := &QryParam{}
	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := u.Message(true, "OK!")
	resp["personal"] = account.GetPerformanceSales("personal", qry.Date, qry.AccountID)
	resp["direct"] = account.GetPerformanceSales("direct", qry.Date, qry.AccountID)
	resp["sub_upline"] = account.GetPerformanceSales("sub_upline", qry.Date, qry.AccountID)
	resp["upline"] = account.GetPerformanceSales("upline", qry.Date, qry.AccountID)
	u.Respond(w, resp)
}

func GetAccessAPI(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	// https://acct.servehpbr.com
	resp := map[string]interface{}{"endpoint": "https://acct.servehpbr.com", "version": 1, "dismissible": true, "download_url": "https://play.blazingsphere.net"}

	u.Respond(w, resp)
}
