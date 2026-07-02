package account

import (
	"fmt"
	"time"

	u "gollux/utils"
	"zerasuite/bookings"
	"zerasuite/shippinglines"
	"zerasuite/yards"
)

type Query struct {
	Table     string `json:"table"`
	Query     string `json:"query"`
	Type      string `json:"type"`
	AccountID int    `json:"account_id"`
}

type AgentSale struct {
	ID          int       `json:"id"`
	AgentID     int       `json:"agent_id"`
	AgentType   string    `json:"agent_type"`
	Sales       float64   `json:"sales"`
	Commission  int       `json:"commission"`
	CoveredDate string    `json:"covered_date"`
	GeneratedAt time.Time `json:"generated_at"`
	Amount      float64   `json:"amount"`
}
type Settlement struct {
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
}
type Setting struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

func (qry *Query) CustomQry() map[string]interface{} {
	var account []Account
	var trans []Transaction
	var admins []Administrator
	var requests []Request
	var yards []yards.Yard
	var shippinglines []shippinglines.ShippingLine
	var settings []Setting
	var bookingslots []bookings.BookingSlot
	var errdb error
	var response map[string]interface{}
	fmt.Println(qry.Query)
	if qry.Table == "accounts" {
		_, errdb = DBM.Query(&account, qry.Query)
	} else if qry.Table == "transactions" {
		_, errdb = DBM.Query(&trans, qry.Query)
	} else if qry.Table == "administrators" {
		_, errdb = DBM.Query(&admins, qry.Query)
	} else if qry.Table == "requests" {
		_, errdb = DBM.Query(&requests, qry.Query)
	} else if qry.Table == "yards" {
		_, errdb = DBM.Query(&yards, qry.Query)
	} else if qry.Table == "shippinglines" {
		_, errdb = DBM.Query(&shippinglines, qry.Query)
	} else if qry.Table == "settings" {
		_, errdb = DBM.Query(&settings, qry.Query)
	} else if qry.Table == "bookingslots" {
		_, errdb = DBM.Query(&bookingslots, qry.Query)
	} else {
		return u.Message(false, "Invalid table!")
	}
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	response = u.Message(true, "Result")
	switch qry.Table {
	case "accounts":
		response["accounts"] = account
	case "transactions":
		response["transactions"] = trans
	case "administrators":
		response["administrators"] = admins
	case "requests":
		response["requests"] = requests
	case "yards":
		response["yards"] = yards
	case "shippinglines":
		response["shippinglines"] = shippinglines
	case "settings":
		response["settings"] = settings
	case "bookingslots":
		response["bookingslots"] = bookingslots
	}

	return response

}

func (qry *Query) AccountUpdate() map[string]interface{} {
	fmt.Println("AccountUpdate")
	res, errdb := DBM.Exec(qry.Query)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	fmt.Println(res.Model())
	if res.RowsAffected() == 0 {
		return u.Message(false, "No rows affected!")
	}
	response := u.Message(true, "Save")
	response[qry.Table] = res
	return response

}

func AccountStats() map[string]interface{} {
	var accounts []Account
	var topplayers []map[string]interface{}
	_, errdb := DBM.Query(&accounts, `select * from accounts order by balance DESC NULLS LAST`)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	credits := 00.00
	verified := 0
	for _, acc := range accounts {
		credits += acc.Balance
		if acc.Remarks == "Verified" {
			if verified < 8 {
				topplayers = append(topplayers, map[string]interface{}{"FirstName": acc.FirstName, "LastName": acc.LastName, "MobileNo": acc.MobileNo, "AccountID": acc.ID, "Balance": acc.Balance})
			}
			verified++
		}

	}
	fmt.Println("verified", verified)
	response := u.Message(true, "Save")
	response["credits"] = credits
	response["verified"] = verified
	response["count"] = len(accounts)
	response["topplayers"] = topplayers
	return response

}

func GenerateAgentSales() {

	var gendata []struct {
		AgentID int     `json:"agent_id"`
		Sales   float64 `json:"sales"`
	}
	yesterday := time.Now().AddDate(0, 0, -1)
	_date := yesterday.Format("2006-01-02")

	if IsGenerateAgentSales(_date) {
		return // already generated
	}

	_, errdb := DBM.Query(&gendata, `select a.agent as agent_id, sum(amount) as sales  from transactions as t left join accounts a on a.id=t.account_id 
where t.type = 'payment' and to_date(t.created_at::text, 'YYYY-MM-DD') = ?   group by a.agent`, _date)
	if errdb != nil {
		panic(errdb)
	}
	for _, data := range gendata {
		var sales AgentSale
		sales.AgentID = data.AgentID
		sales.AgentType = "Agent"
		sales.Sales = data.Sales
		sales.CoveredDate = _date
		sales.GeneratedAt = time.Now()
		errdb = DBM.Insert(&sales)
		if errdb != nil {
			panic(errdb)
		}

	}
}
func IsGenerateAgentSales(_date string) bool {

	var sales []AgentSale
	_, errdb := DBM.Query(&sales, `select * from agent_sales where covered_date=?`, _date)
	if errdb != nil {
		panic(errdb)
	}
	if len(sales) > 0 {
		return true
	}
	return false
}
func IsGenerateAgentSettlement(_date string) bool {

	var sales []Settlement
	_, errdb := DBM.Query(&sales, `select * from settlements where to_date(generated_at::text,'YYYY-MM-DD') = ? `, _date)
	if errdb != nil {
		panic(errdb)
	}
	if len(sales) > 0 {
		return true
	}
	return false
}

func GenerateSettlement() {

	var gendata []struct {
		AgentID    int     `json:"agent_id"`
		Sales      float64 `json:"sales"`
		Commission int     `json:"commission"`
		Amount     float64 `json:"amount"`
	}
	_date := time.Now()
	_today := _date.Format("2006-01-02")
	from_date := time.Now().AddDate(0, 0, -7)
	to_date := time.Now().AddDate(0, 0, -1)
	_to_date := to_date.Format("2006-01-02")
	_from_date := from_date.Format("2006-01-02")

	if IsGenerateAgentSettlement(_today) {
		fmt.Println("Settlement  already generated")
		return // already generated
	}

	_, errdb := DBM.Query(&gendata, `select agent_id, sum(sales) as sales, sum(amount) as amount from agent_sales 
where to_date(covered_date,'YYYY-MM-DD') >= ? and to_date(covered_date,'YYYY-MM-DD') <= ? group by agent_id`, _from_date, _to_date)
	if errdb != nil {
		panic(errdb)
	}
	for _, data := range gendata {
		var sales Settlement
		sales.AgentID = data.AgentID
		sales.AgentType = "Agent"
		sales.Sales = data.Sales
		sales.CoveredDate = _from_date + " to " + _to_date
		sales.GeneratedAt = time.Now()
		sales.Remarks = "Settlement Generated"
		sales.Status = "Generated"
		sales.Amount = data.Amount
		//sales.SetCommission = data.Commission
		errdb = DBM.Insert(&sales)
		if errdb != nil {
			panic(errdb)
		}

	}
}
func PlayersLocation() map[string]interface{} {
	var gendata []Location
	_, errdb := DBM.Query(&gendata, `select * from locations where lat is not null order by id desc limit 500`)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, "Ok")
	response["items"] = gendata
	return response
}

func (qry *Query) AgentSales() map[string]interface{} {
	var sales []AgentSale
	_, errdb := DBM.Query(&sales, qry.Query)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, "Ok")
	response["items"] = sales
	return response
}

func (qry *Query) AgentLocation() map[string]interface{} {
	var locations []Location
	_, errdb := DBM.Query(&locations, qry.Query)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, "Ok")
	response["items"] = locations
	return response
}

func (qry *Query) Settlements() map[string]interface{} {
	var items []Settlement
	_, errdb := DBM.Query(&items, qry.Query)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, "Ok")
	response["items"] = items
	return response
}

func GetPerformanceSales(_type string, _date string, _accountID int) map[string]interface{} {
	fmt.Println("Get Generating agent sales...", _type, _date)
	var gendata []struct {
		AgentID int     `json:"agent_id"`
		Sales   float64 `json:"sales"`
		Comm    int     `json:"comm"`
	}

	qry_str := ""
	comm := 1
	switch _type {
	case "upline":
		comm = UplineCommmission
		qry_str = "select a.upline as agent_id, sum(amount) as sales, ? as comm  from transactions as t left join accounts a on a.id=t.account_id where t.type = 'payment' and to_date(t.created_at::text, 'YYYY-MM-DD') = ? and a.upline=? group by a.upline"
	case "sub_upline":
		comm = SubUplineCommmission
		qry_str = "select a.sub_upline as agent_id, sum(amount) as sales, ? as comm  from transactions as t left join accounts a on a.id=t.account_id where t.type = 'payment' and to_date(t.created_at::text, 'YYYY-MM-DD') = ? and a.sub_upline =? group by a.sub_upline"
	case "direct":
		comm = DirectCommmission
		qry_str = "select a.direct_upline as agent_id, sum(amount) as sales, ? as comm  from transactions as t left join accounts a on a.id=t.account_id where t.type = 'payment' and to_date(t.created_at::text, 'YYYY-MM-DD') = ? and a.direct_upline =? group by a.direct_upline"
	case "personal":
		comm = PersonalCommmission
		qry_str = "select a.id as agent_id, sum(amount) as sales, ? as comm  from transactions as t left join accounts a on a.id=t.account_id where t.type = 'payment' and to_date(t.created_at::text, 'YYYY-MM-DD') = ? and a.id =?  group by a.id"
	default:
		return u.Message(false, _type)
	}

	_, errdb := DBM.Query(&gendata, qry_str, comm, _date, _accountID)
	if errdb != nil {
		panic(errdb)
	}
	response := u.Message(true, _type)
	response["sales"] = gendata
	return response
}
