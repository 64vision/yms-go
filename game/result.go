package game

import (
	"fmt"
	"strconv"
	"strings"

	u "gollux/utils"

	"github.com/dustin/go-humanize"
)

type Rank struct {
	Rank string `json:"rank"`
	Odds string `json:"odds"`
	Ball string `json:"ball"`
	Race string `json:"race"`
}
type Recent struct {
	Raceno string `json:"raceno"`
	B1     string `json:"b1"`
	B2     string `json:"b2"`
	B3     string `json:"b3"`
	B4     string `json:"b4"`
	B5     string `json:"b5"`
	B6     string `json:"b6"`
	B7     string `json:"b7"`
	B8     string `json:"b8"`
	B9     string `json:"b9"`
	B10    string `json:"b10"`
}
type Prize struct {
	Win  string `json:"Win"`
	Game string `json:"game"`
}

func VmixResult() map[string]interface{} {
	var item Race
	var results []Rank
	_, errdb := DBM.Query(&item, `SELECT * from races where status=? order by id desc limit 1`, "Complete") //get last complete
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	parts := strings.Split(item.Result, ",")
	for i, ball := range parts {
		parts[i] = strings.TrimSpace(ball) // Remove spaces
		suffix := "th"
		if (i + 1) == 1 {
			suffix = "st"
		}
		if (i + 1) == 2 {
			suffix = "nd"
		}
		if (i + 1) == 3 {
			suffix = "rd"
		}
		rank := fmt.Sprintf("%d"+suffix, i+1)
		results = append(results, Rank{
			Rank: rank,
			Odds: getOdds(ball, item),
			Ball: ball,
			Race: item.Number,
		})
	}

	response := u.Message(true, "OK")
	response["results"] = results

	return response
}

func getOdds(ball string, race Race) string {
	odds := race.Odds["results"]
	//t.Println(odds)
	if slice, ok := odds.([]interface{}); ok {
		for _, mapball := range slice {
			//fmt.Println("_here", ball["odds"])
			m, ok := mapball.(map[string]interface{})
			if !ok {
				fmt.Println("ball is not a map")

			}
			winball := strconv.Itoa(int(m["ball"].(float64)))
			if winball == ball {
				return fmt.Sprintf("%.2f", m["odds"])
				//fmt.Println("win ball found here 123", winball)
			}
		}

	}
	return "1"
}

func VmixRecentResult() map[string]interface{} {
	var items []Race
	var recents []Recent
	_, errdb := DBM.Query(&items, `SELECT * from races where status=? order by id desc limit 5`, "Complete") //get last complete
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	for _, item := range items {
		parts := strings.Split(item.Result, ",")
		recents = append(recents, Recent{
			Raceno: "RACE NO: " + item.Number,
			B1:     parts[0],
			B2:     parts[1],
			B3:     parts[2],
			B4:     parts[3],
			B5:     parts[4],
			B6:     parts[5],
			B7:     parts[6],
			B8:     parts[7],
			B9:     parts[8],
			B10:    parts[9], //09301261697
		})
	}

	response := u.Message(true, "OK")
	response["results"] = recents

	return response
}

func PrizesVmix() map[string]interface{} {
	var items []Setting
	var prizes []Prize
	_, errdb := DBM.Query(&items, `select * from settings where description = 'Prize' order by id asc`) //get last 09228910124
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	for _, item := range items {
		str := item.Win
		win, err := strconv.ParseFloat(str, 64)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			prizes = append(prizes, Prize{
				Game: item.Game,
				Win:  humanize.CommafWithDigits(win, 2),
			})
		}

	}

	response := u.Message(true, "OK")
	response["results"] = prizes

	return response
}
