package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	u "gollux/utils"
)

type AccountSettlement struct {
	ID            int       `json:"id"`
	AgentID       int       `json:"agent_id"`
	AgentType     string    `json:"agent_type"`
	Sales         float64   `json:"sales"`
	SetCommission int       `json:"set_commission"`
	CoveredDate   string    `json:"covered_date"`
	UpdatedBy     int       `json:"updated_by"`
	GeneratedAt   time.Time `json:"generated_at"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	Remarks       string    `json:"remarks"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	MobileNo      string    `json:"mobile_no"`
}

func GetAccountSettlement(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &AccountSettlement{}

	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := qry.GetAccountSettlements()
	u.Respond(w, resp)
}
func UpdateSettlement(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &AccountSettlement{}

	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	qry.UpdatedBy = r.Context().Value("user").(int)
	resp := qry.Update()
	u.Respond(w, resp)
}
func GetSettlements(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Query{}

	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := QuerySettlements()
	u.Respond(w, resp)
}
func QuerySettlements() map[string]interface{} {
	var items []AccountSettlement
	_, errdb := DBM.Query(&items, `SELECT s.*, a.first_name, a.last_name, a.mobile_no FROM settlements s left join accounts a on a.id=s.agent_id 
order by id desc`)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, "Result")
	response["settlements"] = items
	return response
}
func (settlement *AccountSettlement) Update() map[string]interface{} {

	now := time.Now()
	time := now.Format("2006-01-02 15:04:05")
	settlement.Remarks = settlement.Status + " on " + time

	_, errdb := DBM.Exec(`Update settlements set status=?, remarks=? where id=?`, settlement.Status, settlement.Remarks, settlement.ID)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, "Updated")
	return response
}
func (acct *AccountSettlement) GetAccountSettlements() map[string]interface{} {
	var settlement AccountSettlement
	var sales []AgentSale

	_, errdb := DBM.Query(&settlement, `SELECT s.*, a.first_name, a.last_name, a.mobile_no FROM settlements s left join accounts a on a.id=s.agent_id  where s.id=? order by s.id desc`, acct.ID)
	if errdb != nil {
		panic(errdb)
	}
	_covered_date := strings.Split(settlement.CoveredDate, " to ")

	fmt.Println(_covered_date[0])
	fmt.Println(_covered_date[1])
	_, errdb = DBM.Query(&sales, `select * from agent_sales where  to_date(covered_date,'YYYY-MM-DD') >= ? and to_date(covered_date,'YYYY-MM-DD') <= ? and agent_id=? order by id desc`, _covered_date[0], _covered_date[1], settlement.AgentID)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, "Result")
	response["settlement"] = settlement
	response["sales"] = sales
	return response
}
