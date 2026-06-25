package account

import (
	"fmt"
	"time"
)

const (
	UplineCommmission    = 1 // percent
	SubUplineCommmission = 1 // percent
	DirectCommmission    = 3 // percent
	PersonalCommmission  = 3 // percent
)

func GenerateSales(_type string) {
	fmt.Println("Generating agent sales...", _type)
	var gendata []struct {
		AgentID int     `json:"agent_id"`
		Sales   float64 `json:"sales"`
	}
	yesterday := time.Now().AddDate(0, 0, -1)
	_date := yesterday.Format("2006-01-02")
	// := time.Now().Format("2006-01-02")

	if IsGenerateSales(_date, _type) {
		return // already generated
	}
	qry_str := ""
	comm := 1
	switch _type {
	case "upline":
		comm = UplineCommmission
		qry_str = "select a.upline as agent_id, sum(amount) as sales  from transactions as t left join accounts a on a.id=t.account_id where t.type = 'payment' and to_date(t.created_at::text, 'YYYY-MM-DD') = ? and upline is not null  group by a.upline"
	case "sub_upline":
		comm = SubUplineCommmission
		qry_str = "select a.sub_upline as agent_id, sum(amount) as sales  from transactions as t left join accounts a on a.id=t.account_id where t.type = 'payment' and to_date(t.created_at::text, 'YYYY-MM-DD') = ? and sub_upline is not null  group by a.sub_upline"
	case "direct":
		comm = DirectCommmission
		qry_str = "select a.direct_upline as agent_id, sum(amount) as sales  from transactions as t left join accounts a on a.id=t.account_id where t.type = 'payment' and to_date(t.created_at::text, 'YYYY-MM-DD') = ? and direct_upline is not null  group by a.direct_upline"
	case "personal":
		comm = PersonalCommmission
		qry_str = "select a.id as agent_id, sum(amount) as sales  from transactions as t left join accounts a on a.id=t.account_id where t.type = 'payment' and to_date(t.created_at::text, 'YYYY-MM-DD') = ? and a.id is not null  group by a.id"
	default:
		return
	}

	_, errdb := DBM.Query(&gendata, qry_str, _date)
	if errdb != nil {
		panic(errdb)
	}
	for _, data := range gendata {
		var sales AgentSale
		sales.AgentID = data.AgentID
		sales.AgentType = _type
		sales.Sales = data.Sales
		sales.Commission = comm //(data.Sales * float64(comm)) / 100
		sales.Amount = (data.Sales * float64(comm)) / 100
		sales.CoveredDate = _date
		sales.GeneratedAt = time.Now()
		errdb = DBM.Insert(&sales)
		if errdb != nil {
			panic(errdb)
		}

	}
}

func IsGenerateSales(_date string, _type string) bool {

	var sales []AgentSale
	_, errdb := DBM.Query(&sales, `select * from agent_sales where covered_date=? and agent_type=?`, _date, _type)
	if errdb != nil {
		panic(errdb)
	}
	if len(sales) > 0 {
		return true
	}
	return false
}
