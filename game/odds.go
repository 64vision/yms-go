package game

import (
	"fmt"
	"strconv"

	u "hyperball.com/utils"
)

type Odds struct {
	Ball int     `json:"ball"`
	Odds float64 `json:"odds"`
}

func (race *Race) GetOddsVmix() map[string]interface{} {
	fmt.Println("GetOddsVmix...")
	var items []Odds
	var total struct {
		Amount float64 `json:"amount"`
	}
	//get Settings
	s := &Setting{}
	s.Game = HW
	s = s.Query()
	percentage, err := strconv.ParseFloat(s.Win, 64) // Convert to float64
	if err != nil {
		fmt.Println("GetOddsVmix Error:", err)

	}

	_, err = DBM.Query(&total, `select  sum(amount) as amount from bets where race_id=? and game = ?`, race.ID, HW)
	if err != nil {
		panic(err)
		return u.Message(false, "Failed to create account, connection error")
	}
	less_amount := total.Amount - (total.Amount * (percentage / 100))

	_, errd := DBM.Query(&items, `select ball, trunc((?/total)::NUMERIC(10,3),2) as odds from (select cast(balls as int) as ball, sum(amount) as total from bets where race_id=? and game = ? group by balls) perball order by ball asc`, less_amount, race.ID, HW)
	if errd != nil {
		panic(errd)
		return u.Message(false, "Failed to create account, connection error")
	}
	response := map[string]interface{}{"results": items}
	return response
}
