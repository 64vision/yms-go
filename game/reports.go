package game

import (
	"fmt"
	"strconv"
	"time"

	"gollux/reports"
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
	Table    string `json:"table"`
	Query    string `json:"query"`
	Type     string `json:"type"`
	IntValue int    `json:"int_value"`
}

func GenerateSalesReports(raceID int) {
	DeleteDups(raceID)
	race := CurrentRace(raceID)
	per_game := GenPerGame(raceID, "sales")
	per_agent := GenPerAgent(raceID, "sales")
	per_province := GenPerProvince(raceID, "sales")

	var sales SalesSummary
	sales.RaceID = raceID
	sales.RaceNo = race.Number
	sales.RaceDate = race.Date
	sales.PerGame = per_game
	sales.PerAgent = per_agent
	sales.PerProvince = per_province
	sales.GeneratedAt = time.Now()
	err := reports.DBM.Insert(&sales)
	if err != nil {
		panic(err.Error())
	} else {

		fmt.Println("Per game reports generated")
	}
}

func GenerateWinReports(raceID int) {

	race := CurrentRace(raceID)
	per_game := GenPerGame(raceID, "win")
	per_agent := GenPerAgent(raceID, "win")
	per_province := GenPerProvince(raceID, "win")

	var win WinSummary
	win.RaceID = raceID
	win.RaceNo = race.Number
	win.RaceDate = race.Date
	win.PerGame = per_game
	win.PerAgent = per_agent
	win.PerProvince = per_province
	win.GeneratedAt = time.Now()
	err := reports.DBM.Insert(&win)
	if err != nil {
		panic(err.Error())
	} else {
		go TruncateBets(raceID)
		fmt.Println("Per game reports generated")
	}

}
func TruncateBets(raceID int) {
	time.Sleep(2 * time.Second)
	num := strconv.Itoa(raceID)
	_, err := DBM.Exec(`delete from bets where race_id =` + num)
	if err != nil {
		panic(err.Error())
	}
}
func DeleteDups(raceID int) {
	_, err := reports.DBM.Exec(`delete from sales_summaries where race_id =?`, raceID)
	if err != nil {
		panic(err.Error())
	}
}

func GenPerGame(raceID int, reptype string) []PerGame {
	fmt.Println("GenPerGame", reptype)
	var reports []PerGame
	qryStr := ""
	sales := `select game, sum(amount) as total from bets  where race_id=?
group by  game`
	win := `select game, sum(win) as total from bets  where race_id=? and status='Win'
group by  game`
	if reptype == "sales" {
		qryStr = sales
	} else if reptype == "win" {
		qryStr = win
	}

	_, err := DBM.Query(&reports, qryStr, raceID)
	if err != nil {
		panic(err.Error())
		return nil
	}
	fmt.Println("GenPerGame reports", raceID)
	return reports
}

func GenPerAgent(raceID int, reptype string) []PerAgent {
	var reports []PerAgent

	qryStr := ""
	sales := `select agent, province, game,  sum(amount) as total from bets  where race_id=?
group by  agent, province, game`
	win := `select agent, province, game,  sum(win) as total from bets  where race_id=? and status='Win'
group by  agent, province, game`
	if reptype == "sales" {
		qryStr = sales
	} else if reptype == "win" {
		qryStr = win
	}

	_, err := DBM.Query(&reports, qryStr, raceID)
	if err != nil {
		panic(err.Error())
		return nil
	}
	return reports
}

func GenPerProvince(raceID int, reptype string) []PerProvince {
	var reports []PerProvince

	qryStr := ""
	sales := `select province, city, game,  sum(amount) as total from bets  where race_id=?
group by  province, city, game`
	win := `select province, city, game,  sum(win) as total from bets  where race_id=? and status='Win'
group by  province, city, game`
	if reptype == "sales" {
		qryStr = sales
	} else if reptype == "win" {
		qryStr = win
	}

	_, err := DBM.Query(&reports, qryStr, raceID)
	if err != nil {
		panic(err.Error())
		return nil
	}
	return reports
}

func (qry *Query) QryBetLogs() map[string]interface{} {
	var bets []BetLog
	var race []Race
	var err error
	if qry.Type == "bets" {
		_, err = DBM.Query(&bets, qry.Query)
	} else if qry.Type == "race results" {
		_, err = DBM.Query(&race, qry.Query)
	}

	if err != nil {
		panic(err.Error())
		return nil
	}
	response := u.Message(true, "Result")
	if qry.Type == "bets" {
		response["bets"] = bets
	} else if qry.Type == "race results" {
		response["race"] = race
	}

	return response
}
