package reports

import (
	"encoding/json"
	"net/http"
	"time"

	u "gollux/utils"
)

type PerGame struct {
	Game  string  `json:"game"`
	Total float64 `json:"total"`
}

type PerAgent struct {
	Agent    int     `json:"agent"`
	Province string  `json:"province"`
	Game     string  `json:"game"`
	Total    float64 `json:"total"`
}
type PerProvince struct {
	Province string  `json:"province"`
	City     string  `json:"city"`
	Game     string  `json:"game"`
	Total    float64 `json:"total"`
}
type SalesSummary struct {
	ID          int           `json:"id"`
	RaceID      int           `json:"race_id"`
	RaceNo      string        `json:"race_no"`
	RaceDate    string        `json:"race_date"`
	PerGame     []PerGame     `pg:",array" json:"per_game"`
	PerAgent    []PerAgent    `pg:",array" json:"per_agent"`
	PerProvince []PerProvince `pg:",array" json:"per_province"`
	GeneratedAt time.Time     `json:"generated_at"`
}
type WinSummary struct {
	ID          int           `json:"id"`
	RaceID      int           `json:"race_id"`
	RaceNo      string        `json:"race_no"`
	RaceDate    string        `json:"race_date"`
	PerGame     []PerGame     `pg:",array" json:"per_game"`
	PerAgent    []PerAgent    `pg:",array" json:"per_agent"`
	PerProvince []PerProvince `pg:",array" json:"per_province"`
	GeneratedAt time.Time     `json:"generated_at"`
}
type Query struct {
	Table string `json:"table"`
	Query string `json:"query"`
	Type  string `json:"type"`
}

func GetReports(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Query{}
	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := qry.QueryReports()

	u.Respond(w, resp)
}

func (qry *Query) QueryReports() map[string]interface{} {
	var sales []SalesSummary
	var wins []WinSummary
	var errdb error
	var response map[string]interface{}

	if qry.Table == "sales_summaries" {
		_, errdb = DBM.Query(&sales, qry.Query)
	} else if qry.Table == "win_summaries" {
		_, errdb = DBM.Query(&wins, qry.Query)
	} else {
		return u.Message(false, "Invalid table!")
	}
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	response = u.Message(true, "Result")
	switch qry.Table {
	case "sales_summaries":
		response["sales"] = sales
	case "win_summaries":
		response["win"] = wins
	}

	return response
}

func Performance() map[string]interface{} {
	var gendata struct {
		RaceDate string      `json:"race_date"`
		Bets     interface{} `json:"bets"`
	}
	_, errdb := DBM.Query(&gendata, `
SELECT race_date,  JSON_AGG(
        JSON_BUILD_OBJECT(
            'sales', per_game,
			   'number', race_no
        )
    ) AS bets from sales_summaries where  to_date(race_date, 'YYYY-MM-DD') >= NOW() - INTERVAL '30 days' group by race_date 
	order by race_date asc
`)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	response := u.Message(true, "Result")
	response["performance"] = gendata
	return response
}
