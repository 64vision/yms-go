package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"hyperball.com/account"
	u "hyperball.com/utils"
)

func GetSettings(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Setting{}
	sys := &SysSetting{}
	resp := u.Message(false, "Ok")
	resp["games"] = qry.Get()
	resp["sys"] = sys.Get()
	u.Respond(w, resp)
}
func StreamSettings(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	sys := &SysSetting{}
	resp := u.Message(false, "Ok")
	resp["sys"] = sys.Get()
	u.Respond(w, resp)
}

func PlaceHyperWinBet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PlaceHyperWinBet")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	bet := &BetWin{}
	total := 00.00
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(bet) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	race := CurrentRace(bet.RaceID)
	if bet.Balls == "" {
		u.Respond(w, u.Message(false, "Select racing ball!"))
		return
	}
	if race.Status != "Open" {
		u.Respond(w, u.Message(false, "Betting is already closed for this race!"))
		return
	}
	if race.Number != bet.RaceNo {
		u.Respond(w, u.Message(false, "Race Number mismatch!"))
		return
	}
	if bet.RaceDate != race.Date {
		u.Respond(w, u.Message(false, "Race date mismatch!"))
		return
	}
	if bet.AccountID != r.Context().Value("user").(int) {
		u.Respond(w, u.Message(false, "Account ID mismatch!"))
		return
	}
	acct := account.InquireAccount(bet.AccountID)

	amounts := strings.Split(bet.Amount, ",")
	for _, amount := range amounts {

		amt, err := strconv.ParseFloat(amount, 64)

		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("amt:", amt)
			total += amt
		}

	}

	if acct.Balance < total {
		u.Respond(w, u.Message(false, "Insufficient funds!"))
		return
	}

	// validresp := bet.Validate()
	// if !validresp["status"].(bool) {
	// 	u.Respond(w, validresp)
	// 	return
	// }
	resp = bet.PlaceHyperWin()

	u.Respond(w, resp)
}

func PlaceBet(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	bet := &Bet{}
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(bet) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	race := CurrentRace(bet.RaceID)

	if race.Status != "Open" {
		u.Respond(w, u.Message(false, "Betting is already closed for this race!"))
		return
	}
	if race.Number != bet.RaceNo {
		u.Respond(w, u.Message(false, "Race Number mismatch!"))
		return
	}
	if bet.RaceDate != race.Date {
		u.Respond(w, u.Message(false, "Race date mismatch!"))
		return
	}
	if bet.AccountID != r.Context().Value("user").(int) {
		u.Respond(w, u.Message(false, "Account ID mismatch!"))
		return
	}
	if bet.BetType == RBL {
		if bet.Amount < 30 {
			u.Respond(w, u.Message(false, "30 Minimum Bet is required for Shuffle Mode!"))
			return
		}
	}
	acct := account.InquireAccount(bet.AccountID)
	if acct.Balance < bet.Amount {
		u.Respond(w, u.Message(false, "Insufficient funds!"))
		return
	}

	validresp := bet.Validate()
	if !validresp["status"].(bool) {
		u.Respond(w, validresp)
		return
	}
	if bet.Game == HA {
		if bet.Amount < HYPER_TOTALBET {
			u.Respond(w, u.Message(false, "Invalid hyper all bet amount!"))
			return
		}
		resp = bet.PlaceHyperAll()
	} else {
		resp = bet.Place()
	}

	u.Respond(w, resp)
}

func OpenRace(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	race := &Race{}
	err := json.NewDecoder(r.Body).Decode(race) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if race.CheckOpen() != 0 {
		u.Respond(w, u.Message(false, "Race is already open!"))
		return
	}
	race.Number = race.GetRaceNumber()

	race.CreatedBy = r.Context().Value("user").(int)
	resp := race.Open()
	u.Respond(w, resp)
}

func KeyInRace(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	race := &Race{}
	err := json.NewDecoder(r.Body).Decode(race) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	fmt.Println(race)
	resp := race.KeyIn()
	u.Respond(w, resp)
}

func UpdateRaceLap(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	race := &Race{}
	err := json.NewDecoder(r.Body).Decode(race) //decode the request body into struct and failed if any error occur
	if err != nil {
		panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := race.UpdateLap()
	u.Respond(w, resp)
}

func UpdateRace(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	race := &Race{}
	err := json.NewDecoder(r.Body).Decode(race) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := race.Update()
	u.Respond(w, resp)
}
func GetResultVmix(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	resp := VmixResult()
	u.Respond(w, resp)
}
func GetPrizesVmix(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	resp := PrizesVmix()
	u.Respond(w, resp)
}
func GetRecentResultVmix(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	resp := VmixRecentResult()
	u.Respond(w, resp)
}
func GetLapVmix(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	resp := CurrentRaceLap()
	u.Respond(w, resp)
}

func GetOddsVmix(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	race := &Race{}
	race.ID = race.CurrentRace()
	resp := race.GetOddsVmix()
	u.Respond(w, resp)
}

func GetOpenRace(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	race := &Race{}

	resp := race.GetOpen()
	u.Respond(w, resp)
}

func GetRaceList(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	race := &QryParams{}
	err := json.NewDecoder(r.Body).Decode(race) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := race.QryList()
	u.Respond(w, resp)
}

func GetHyperAllNum(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	resp := u.Message(true, "Ok")
	resp["balls"] = u.GenerateUnique10DigitNumber()
	u.Respond(w, resp)
}

func GetWager(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	bet := &Bet{}
	err := json.NewDecoder(r.Body).Decode(bet) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := bet.Wager()
	u.Respond(w, resp)
}

func QueryBets(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Query{}
	resp := u.Message(true, "Ok")
	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if qry.Type == "per game" {
		pergame := GenPerGame(qry.IntValue, qry.Table)
		resp["pergame"] = pergame
	}

	u.Respond(w, resp)
}

func GetBetLogs(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	qry := &Query{}
	err := json.NewDecoder(r.Body).Decode(qry) //decode the request body into struct and failed if any error occur
	if err != nil {
		//panic(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := qry.QryBetLogs()
	u.Respond(w, resp)
}
