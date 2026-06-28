package game

import (
	"fmt"
	"strconv"
	"strings"

	"gollux/account"
	"gollux/sms"
)

func (race *Race) ParseWinner() {
	// parse the result to find the winner
	// and set the winner id and name in the race
	// ...
	balls := strings.Split(race.Result, ",")

	HyperWin(balls, race.ID)
	Hyper2(balls, race.ID)
	Hyper3(balls, race.ID)
	//Hyper4(balls, race.ID)
	Hyper5(balls, race.ID)
	Super6(balls, race.ID)
	Super7(balls, race.ID)
	Super8(balls, race.ID)
	Super9(balls, race.ID)

	//go GenerateReports(raceID int)
}

func QryBets(raceID int, game string) []Bet {
	var bets []Bet
	_, err := DBM.Query(&bets, `SELECT * from bets where race_id=? and game=?`, raceID, game)
	if err != nil {
		panic(err)
	}
	return bets
}

func QrySuperWin(raceID int, game string, winballs string) []Bet {
	var bets []Bet
	_, err := DBM.Query(&bets, `SELECT * from bets where race_id=? and game=? and balls=?`, raceID, game, winballs)
	if err != nil {
		panic(err)
	}
	return bets
}
func Hyper3(balls []string, raceID int) {
	var s Setting
	s.Game = H3
	s = *s.Query()
	winstr := s.Win
	winperunit, err := strconv.ParseFloat(winstr, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	winballs := balls[0] + "," + balls[1] + "," + balls[2]
	bets := QryBets(raceID, H3)
	fmt.Println(bets)
	for _, bet := range bets {
		if bet.BetType == "RBL" {
			fmt.Println("RBL: ", bet.Balls, winballs)
			bet.ExtractShuffle(winballs, winperunit)
		} else {
			fmt.Println("STD: ", bet.Balls, winballs)
			fmt.Println("Bet: ", bet)
			_, err = DBM.Exec(`update bets set win=amount*?, status=? where id=? and race_id=? and balls=? and game=?`, winperunit, "Win", bet.ID, raceID, winballs, H3)
			if err != nil {
				panic(err)
			}
		}

	}

	go LoadWins(raceID, H3)
}
func Super9(balls []string, raceID int) {
	var s Setting
	s.Game = S9
	s = *s.Query()
	winstr := s.Win
	winperunit, err := strconv.ParseFloat(winstr, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	winballs := balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4] + "," + balls[5] + "," + balls[6] + "," + balls[7] + "," + balls[8]
	winbets := QrySuperWin(raceID, S9, winballs)
	if len(winbets) > 0 {
		_, err = DBM.Exec(`update bets set win=?, status=? where race_id=? and balls=? and game=?`, winperunit/float64(len(winbets)), "Win", raceID, winballs, S9)
		if err != nil {
			panic(err)
		}
		go LoadWins(raceID, S9)
	}
}
func Super8(balls []string, raceID int) {
	var s Setting
	s.Game = S8
	s = *s.Query()
	winstr := s.Win
	winperunit, err := strconv.ParseFloat(winstr, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	winballs := balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4] + "," + balls[5] + "," + balls[6] + "," + balls[7]
	winbets := QrySuperWin(raceID, S8, winballs)
	if len(winbets) > 0 {
		_, err = DBM.Exec(`update bets set win=?, status=? where race_id=? and balls=? and game=?`, winperunit/float64(len(winbets)), "Win", raceID, winballs, S8)
		if err != nil {
			panic(err)
		}
		go LoadWins(raceID, S8)
	}
}
func Super7(balls []string, raceID int) {
	var s Setting
	s.Game = S7
	s = *s.Query()
	winstr := s.Win
	winperunit, err := strconv.ParseFloat(winstr, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	winballs := balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4] + "," + balls[5] + "," + balls[6]
	winbets := QrySuperWin(raceID, S7, winballs)
	if len(winbets) > 0 {
		_, err = DBM.Exec(`update bets set win=?, status=? where race_id=? and balls=? and game=?`, winperunit/float64(len(winbets)), "Win", raceID, winballs, S7)
		if err != nil {
			panic(err)
		}
		go LoadWins(raceID, S7)
	}
}
func Super6(balls []string, raceID int) {
	var s Setting

	s.Game = S6
	s = *s.Query()
	winstr := s.Win
	winperunit, err := strconv.ParseFloat(winstr, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	winballs := balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4] + "," + balls[5]
	winbets := QrySuperWin(raceID, S6, winballs)
	if len(winbets) > 0 {
		_, err = DBM.Exec(`update bets set win=?, status=? where race_id=? and balls=? and game=?`, winperunit/float64(len(winbets)), "Win", raceID, winballs, S6)
		if err != nil {
			panic(err)
		}
		go LoadWins(raceID, S6)
	}

}
func Hyper5(balls []string, raceID int) {
	var s Setting
	s.Game = H5
	s = *s.Query()
	winstr := s.Win
	winperunit, err := strconv.ParseFloat(winstr, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	winballs := balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4]
	_, err = DBM.Exec(`update bets set win=amount*?, status=? where race_id=? and balls=? and game=?`, winperunit, "Win", raceID, winballs, H5)
	if err != nil {
		panic(err)
	}
	go LoadWins(raceID, H5)
}
func Hyper4(balls []string, raceID int) {
	var s Setting
	s.Game = H4
	s = *s.Query()
	winstr := s.Win
	winperunit, err := strconv.ParseFloat(winstr, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	winballs := balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3]
	_, err = DBM.Exec(`update bets set win=amount*?, status=? where race_id=? and balls=? and game=?`, winperunit, "Win", raceID, winballs, H4)
	if err != nil {
		panic(err)
	}
	go LoadWins(raceID, H4)
}

func Hyper2(balls []string, raceID int) {
	var s Setting
	s.Game = H2
	s = *s.Query()
	winstr := s.Win
	winperunit, err := strconv.ParseFloat(winstr, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	winballs := balls[0] + "," + balls[1]
	_, err = DBM.Exec(`update bets set win=amount*?, status=? where race_id=? and balls=? and game=?`, winperunit, "Win", raceID, winballs, H2)
	if err != nil {
		panic(err)
	}
	go LoadWins(raceID, H2)
}
func HyperWin(balls []string, raceID int) {
	race := CurrentRace(raceID)
	odds := race.Odds["results"]
	odds_miltiplier := 0.0
	fmt.Println("win ball:", balls[0])
	fmt.Println("race odds:", odds)

	if slice, ok := odds.([]interface{}); ok {
		//	fmt.Println(slice)
		for _, ball := range slice {
			//fmt.Println("_here", ball["odds"])
			m, ok := ball.(map[string]interface{})
			if !ok {
				fmt.Println("ball is not a map")
				return
			}
			winball := strconv.Itoa(int(m["ball"].(float64)))
			//fmt.Println("win ball found here 123", winball)
			if winball == balls[0] {

				odds_miltiplier = m["odds"].(float64)
			}
		}
		fmt.Println("win ball multiplier", odds_miltiplier)
		_, err := DBM.Exec(`update bets set win=amount*?, status=? where race_id=? and balls=? and game=?`, odds_miltiplier, "Win", raceID, balls[0], HW)
		if err != nil {
			panic(err)
		}
		go LoadWins(raceID, HW)

	} else {
		fmt.Println("No bet on HW")
		//panic("Odds was not saved!")
	}

}

func LoadWins(raceID int, game string) {
	var bets []Bet
	_, err := DBM.Query(&bets, `select * from bets where status=? and race_id=? and game=?`, "Win", raceID, game)
	if err != nil {
		panic(err)
	}
	for _, bet := range bets {
		var trans account.Transaction
		trans.AccountID = bet.AccountID
		trans.Description = "Winning bet"
		trans.Type = account.Deposit
		trans.Amount = bet.Win
		trans.RefNo = bet.ID
		trans.SenderID = 8

		details := map[string]interface{}{"type": trans.Description}
		details["info"] = bet
		trans.Details = details
		wresp := trans.Add()
		if !wresp["status"].(bool) {
			panic(wresp["message"].(string))
		} else {
			//Trigger a notification
			go sms.SendWinNotification(trans.AccountID)
		}

	}
}
