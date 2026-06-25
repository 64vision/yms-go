package game

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"hyperball.com/account"
	u "hyperball.com/utils"
)

type Bet struct {
	ID         int       `json:"id"`
	AccountID  int       `json:"account_id"`
	RaceID     int       `json:"race_id"`
	RaceDate   string    `json:"race_date"`
	RaceNo     string    `json:"race_no"`
	Game       string    `json:"game"`
	Balls      string    `json:"balls"`
	BetType    string    `json:"bet_type"`
	Code       string    `json:"code"`
	Amount     float64   `json:"amount"`
	Win        float64   `json:"win"`
	Status     string    `json:"status"`
	TransDt    time.Time `json:"trans_dt"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	Region     string    `json:"region"`
	Agent      int       `json:"agent"`
	Commission int       `json:"Commission"`
}

type BetLog struct {
	ID         int       `json:"id"`
	AccountID  int       `json:"account_id"`
	RaceID     int       `json:"race_id"`
	RaceDate   string    `json:"race_date"`
	RaceNo     string    `json:"race_no"`
	Game       string    `json:"game"`
	Balls      string    `json:"balls"`
	BetType    string    `json:"bet_type"`
	Code       string    `json:"code"`
	Amount     float64   `json:"amount"`
	Win        float64   `json:"win"`
	Status     string    `json:"status"`
	TransDt    time.Time `json:"trans_dt"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	Region     string    `json:"region"`
	Agent      int       `json:"agent"`
	Commission int       `json:"commission"`
}

type BetWin struct {
	ID         int    `json:"id"`
	AccountID  int    `json:"account_id"`
	RaceID     int    `json:"race_id"`
	RaceDate   string `json:"race_date"`
	RaceNo     string `json:"race_no"`
	BetType    string `json:"bet_type"`
	Game       string `json:"game"`
	Balls      string `json:"balls"`
	Amount     string `json:"amount"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Agent      int    `json:"agent"`
	Commission int    `json:"Commission"`
}

type Wage struct {
	Game string        `json:"game"`
	Bets []interface{} `json:"bets"`
}

const (
	HW             = "HW"  //"Hyper Win"
	H2             = "H2"  //"Hyper 2"
	H3             = "H3"  //"Hyper 3"
	H4             = "H4"  //"Hyper 4"
	H5             = "H5"  //"Hyper 5"
	S6             = "S6"  //"Super 6"
	S7             = "S7"  //"Super 7"
	S8             = "S8"  //"Super 8"
	S9             = "S9"  //"Super 9"
	HA             = "HA"  //"Hyper All"
	STD            = "STD" //" Standard"
	RBL            = "RBL" //"Shuffle
	HYPER_MINBET   = 5     //"
	SUPER_MINBET   = 20    //"
	HYPER_TOTALBET = 100   //"
)

func (bet *Bet) Log() {
	var log BetLog
	log.ID = bet.ID
	log.AccountID = bet.AccountID
	log.RaceID = bet.RaceID
	log.RaceDate = bet.RaceDate
	log.RaceNo = bet.RaceNo
	log.BetType = bet.BetType
	log.Province = bet.Province
	log.City = bet.City
	log.Agent = bet.Agent
	log.Commission = bet.Commission
	log.TransDt = bet.TransDt
	log.Game = bet.Game
	log.Status = bet.Status
	log.Code = bet.Code
	log.Balls = bet.Balls
	log.Amount = bet.Amount
	errdb := DBM.Insert(&log)
	if errdb != nil {
		panic(errdb)

	} else {
		fmt.Println("log.Inserted")

	}

}
func (win *BetWin) PlaceHyperWin() map[string]interface{} {

	win.Balls = strings.ReplaceAll(win.Balls, " ", "")
	win.Amount = strings.ReplaceAll(win.Amount, " ", "")
	balls := strings.Split(win.Balls, ",")
	amounts := strings.Split(win.Amount, ",")
	if len(balls) != len(amounts) {
		return u.Message(false, "Number of balls and amounts does not match")
	}

	for index, ball := range balls {
		var bet Bet

		bet.AccountID = win.AccountID
		bet.RaceID = win.RaceID
		bet.RaceDate = win.RaceDate
		bet.RaceNo = win.RaceNo
		bet.BetType = win.BetType
		bet.Province = win.Province
		bet.City = win.City
		bet.Agent = win.Agent
		bet.Commission = win.Commission

		bet.TransDt = time.Now()
		bet.Game = HW
		bet.Status = "Bet"
		bet.Code = bet.Game
		bet.Balls = ball
		amt, err := strconv.ParseFloat(amounts[index], 64)

		if err != nil {
			fmt.Println("Error:", err)
			//panic(err)
		} else {

			bet.Amount = amt
			//fmt.Println("bet.Amount:", amt)
		}

		errdb := DBM.Insert(&bet)
		if errdb != nil {
			//	panic(errdb)
			return u.Message(false, errdb.Error())
		} else {
			//log bet
			bet.Log()

		}

		//set transaction
		var trans account.Transaction
		trans.AccountID = win.AccountID
		trans.Type = account.Payment
		trans.Amount = bet.Amount
		trans.Description = "Bet"
		trans.SenderID = 10
		trans.RefNo = bet.ID

		details := map[string]interface{}{"type": trans.Description}
		details["info"] = bet

		trans.Details = details

		wresp := trans.Add()
		if !wresp["status"].(bool) {
			return u.Message(false, wresp["message"].(string))
		}
	} //loop end here

	//bet.Combination() // save combination
	acct := account.InquireAccount(win.AccountID)
	response := u.Message(true, "Bet posted successfully!")
	//response["bet"] = bet
	response["account"] = acct
	return response
}
func (bet *Bet) PlaceHyperAll() map[string]interface{} {
	games := []string{HW, H2, H3, H5, S6, S7, S8, S9}
	bet.Balls = strings.ReplaceAll(bet.Balls, " ", "")
	balls := strings.Split(bet.Balls, ",")
	for _, game := range games {
		bet.TransDt = time.Now()
		bet.Game = game
		bet.Status = "Bet"
		bet.Code = HA // u.GenCharCode(6)

		if game == HW || game == H2 || game == H3 || game == H4 || game == H5 {
			bet.Amount = HYPER_MINBET
		} else {
			bet.Amount = SUPER_MINBET
		}

		switch game {
		case HW:
			bet.Balls = balls[0]
		case H2:
			bet.Balls = balls[0] + "," + balls[1]
		case H3:
			bet.Balls = balls[0] + "," + balls[1] + "," + balls[2]
		case H4:
			bet.Balls = balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3]
		case H5:
			bet.Balls = balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4]
		case S6:
			bet.Balls = balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4] + "," + balls[5]
		case S7:
			bet.Balls = balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4] + "," + balls[5] + "," + balls[6]
		case S8:
			bet.Balls = balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4] + "," + balls[5] + "," + balls[6] + "," + balls[7]
		case S9:
			bet.Balls = balls[0] + "," + balls[1] + "," + balls[2] + "," + balls[3] + "," + balls[4] + "," + balls[5] + "," + balls[6] + "," + balls[7] + "," + balls[8]
		default:
			bet.Balls = balls[0]
		}

		errdb := DBM.Insert(bet)
		if errdb != nil {
			panic(errdb)
			return u.Message(false, errdb.Error())
		} else {
			if game == S9 {
				bet.Game = HA
				bet.Amount = HYPER_TOTALBET
				bet.Log()
			}
		}

		//set transaction
		var trans account.Transaction
		trans.AccountID = bet.AccountID
		trans.Type = account.Payment
		trans.Amount = bet.Amount
		trans.Description = "Bet"
		trans.SenderID = 10
		trans.RefNo = bet.ID

		details := map[string]interface{}{"type": trans.Description}
		details["info"] = bet
		trans.Details = details
		wresp := trans.Add()
		if !wresp["status"].(bool) {
			return u.Message(false, wresp["message"].(string))
		}
	} //loop end here
	//bet.Combination() // save combination
	acct := account.InquireAccount(bet.AccountID)
	response := u.Message(true, "Bet posted successfully!")
	//response["bet"] = bet
	response["account"] = acct
	return response
}

func (bet *Bet) Place() map[string]interface{} {
	bet.TransDt = time.Now()
	bet.Status = "Bet"
	bet.Code = bet.Game
	bet.Balls = strings.ReplaceAll(bet.Balls, " ", "")
	errdb := DBM.Insert(bet)
	if errdb != nil {
		panic(errdb)
		return u.Message(false, errdb.Error())
	} else {
		bet.Log()

	}

	//set transaction
	var trans account.Transaction
	trans.AccountID = bet.AccountID
	trans.Type = account.Payment
	trans.Amount = bet.Amount
	trans.Description = "Bet"
	trans.SenderID = 10
	trans.RefNo = bet.ID

	details := map[string]interface{}{"type": trans.Description}
	details["info"] = bet

	trans.Details = details

	wresp := trans.Add()
	if !wresp["status"].(bool) {
		return u.Message(false, wresp["message"].(string))
	}

	//bet.Combination() // save combination
	acct := account.InquireAccount(bet.AccountID)
	response := u.Message(true, "Bet posted successfully!")
	//response["bet"] = bet
	response["account"] = acct
	return response
}

func (bet *Bet) Validate() map[string]interface{} {

	bet.Balls = strings.ReplaceAll(bet.Balls, " ", "")
	balls := strings.Split(bet.Balls, ",")

	fmt.Println("Game", bet.Game)
	fmt.Println("ball len:", len(balls))

	if bet.Game != HA && bet.Game != HW && bet.Game != H2 && bet.Game != H3 && bet.Game != H4 && bet.Game != H5 && bet.Game != S6 && bet.Game != S7 && bet.Game != S8 && bet.Game != S9 {
		return u.Message(false, "Please select game!")
	}

	if bet.Balls == "" {
		return u.Message(false, "Please select balls!")
	}
	if bet.Game == "" {
		return u.Message(false, "Please select game!")
	}
	if bet.Game == HW {
		if len(balls) != 1 {
			return u.Message(false, "Please select 1 ball!")
		}
	}
	if bet.Game == H2 {
		if len(balls) != 2 {
			return u.Message(false, "Please select 2 racing balls!")
		}
	}
	if bet.Game == H3 {
		if len(balls) != 3 {
			return u.Message(false, "Please select 3 racing balls!")
		}
	}
	if bet.Game == H4 {
		if len(balls) != 4 {
			return u.Message(false, "Please select  4 racing balls!")
		}
	}
	if bet.Game == H5 {
		if len(balls) != 5 {
			return u.Message(false, "Please select  5 racing balls!")
		}
	}
	if bet.Game == S6 {
		if len(balls) != 6 {
			return u.Message(false, "Please select  6 racing balls!")
		}
	}
	if bet.Game == S7 {
		if len(balls) != 7 {
			return u.Message(false, "Please select  7 racing balls!")
		}
	}
	if bet.Game == S8 {
		if len(balls) != 8 {
			return u.Message(false, "Please select  8 racing balls!")
		}
	}
	if bet.Game == S9 || bet.Game == HA {
		fmt.Println("ball len p1:", len(balls))
		if len(balls) != 9 {

			return u.Message(false, "Please select  9 racing balls!")
		}
	}

	if u.HasDuplicates(balls) {
		return u.Message(false, "Check duplicates racing balls!")
	}

	return u.Message(true, "OK")
}

func (bet *Bet) Wager() map[string]interface{} {
	var wage []Wage
	_, errdb := DBM.Query(&wage, `SELECT game,  JSON_AGG(
        JSON_BUILD_OBJECT(
			'ball', balls,
            'amount', amount,
			'bet_type', bet_type
        )
    ) AS bets from bet_logs where account_id=? and race_id=?  group by game`, bet.AccountID, bet.RaceID)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	response := u.Message(true, "Result")
	//response["bet"] = bet
	response["wager"] = wage
	return response
}
