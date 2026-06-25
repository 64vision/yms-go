package game

import (
	"strconv"
	"strings"

	u "hyperball.com/utils"
)

type Combination struct {
	ID        int     `json:"ID"`
	AccountID int     `json:"account_id"`
	BetID     int     `json:"bet_id"`
	RaceID    int     `json:"race_id"`
	Game      string  `json:"game"`
	BetType   string  `json:"bet_type"`
	Amount    float64 `json:"amount"`
	Balls     string  `json:"balls"`
	Status    string  `json:"status"`
}

func (bet *Bet) ExtractShuffle(resultballs string, winperunit float64) map[string]interface{} {
	balls := strings.Split(bet.Balls, ",")
	ball1, _ := strconv.Atoi(balls[0])
	ball2, _ := strconv.Atoi(balls[1])
	ball3, _ := strconv.Atoi(balls[2])
	amount := bet.Amount / 6

	arr := []int{ball1, ball2, ball3}
	resArr := u.Permutations(arr)
	for _, combi := range resArr {
		if genCombi(combi) == resultballs {
			_, err := DBM.Exec(`update bets set win=?, status=? where id=? and race_id=? and game=?`, amount*winperunit, "Win", bet.ID, bet.RaceID, H3)
			if err != nil {
				panic(err)
			}
			return u.Message(true, "Ok")
		}

	}
	return u.Message(true, "Ok")
}

func genCombi(combi []int) string {
	balls := make([]string, len(combi))
	for i, ball := range combi {
		balls[i] = strconv.Itoa(ball)
	}
	return strings.Join(balls, ",")

}
