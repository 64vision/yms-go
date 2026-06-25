package game

import (
	"fmt"
	"time"

	"hyperball.com/sms"
	u "hyperball.com/utils"
)

type Race struct {
	ID        int                    `json:"id"`
	Date      string                 `json:"date"`
	CreatedAt time.Time              `json:"created_at"`
	CreatedBy int                    `json:"created_by"`
	Result    string                 `json:"result"`
	Completed time.Time              `json:"completed"`
	Status    string                 `json:"status"`
	Remarks   string                 `json:"remarks"`
	Number    string                 `json:"number"`
	Settings  map[string]interface{} `json:"settings"`
	Odds      map[string]interface{} `json:"odds"`
	Lap       map[string]interface{} `json:"lap"`
}
type QryParams struct {
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}

func (entry *Race) Open() map[string]interface{} {
	//clear noti
	go sms.ClearWinNotification()
	var race Race
	var s Setting
	now := time.Now()
	lap := map[string]interface{}{
		"date": now.Format("Monday, Jan 2, 2006"),
		"race": entry.Number,
		"lap":  "0/1",
	}
	race.CreatedAt = time.Now()
	race.CreatedBy = entry.CreatedBy
	race.Status = "Open"
	race.Remarks = "Open for betting"
	race.Settings = s.Get()
	race.Date = race.CreatedAt.Format("2006-01-02")
	race.Number = entry.Number
	race.Lap = lap
	errdb := DBM.Insert(&race)
	if errdb != nil {
		//	panic(errdb)
		return u.Message(false, errdb.Error())
	}
	response := u.Message(true, "New race opened successfully!")
	response["race"] = &race
	return response
}

func (qry *QryParams) QryList() map[string]interface{} {
	var err error
	var items []Race
	_, err = DBM.Query(&items, `select * from races where to_date(date,'YYYY-MM-DD') >= to_date(?,'YYYY-MM-DD') 
	and  to_date(date,'YYYY-MM-DD') <= to_date(?,'YYYY-MM-DD')
	order by id desc`, qry.FromDate, qry.ToDate) //status to date
	if err != nil {
		panic(err)
		return u.Message(false, "Failed to create account, connection error")
	}
	response := u.Message(true, "Results")
	response["items"] = items
	return response
}

func (entry *Race) GetRaceNumber() string {
	current := time.Now()
	dateNow := current.Format("2006-01-02")
	var race []Race
	res, err := DBM.Query(&race, `SELECT * from races where date=?`, dateNow)
	if err != nil {
		panic(err)
		return "0"
	}
	raceCount := res.RowsReturned() + 1
	return fmt.Sprintf("%d", raceCount)
}

func (entry *Race) CheckOpen() int {
	var race Race
	_, err := DBM.Query(&race, `SELECT * from races where status=? OR status=?`, "Open", "Closed")
	if err != nil {
		panic(err)
		return 0
	}
	return race.ID
}

func (entry *Race) CurrentRace() int {
	var race Race
	_, err := DBM.Query(&race, `SELECT * from races where status=? OR status=?`, "Open", "Closed")
	if err != nil {
		panic(err)
		return 0
	}
	return race.ID
}

func CurrentRaceLap() map[string]interface{} {
	var race Race
	_, err := DBM.Query(&race, `SELECT * from races where status=? OR status=?`, "Open", "Closed")
	if err != nil {
		panic(err)
		//return 0
	}
	return race.Lap
}
func CurrentRace(raceID int) *Race {
	var race Race
	_, err := DBM.Query(&race, `SELECT * from races where id=?`, raceID)
	if err != nil {
		panic(err)
		return nil
	}
	return &race
}
func (entry *Race) GetOpen() map[string]interface{} {
	var race Race
	res, err := DBM.Query(&race, `SELECT * from races where status=? or status=?`, "Open", "Closed")
	if err != nil {
		panic(err)
		return u.Message(false, "Error querying open race")
	}
	if res.RowsReturned() == 0 {
		return u.Message(false, "No found open race")
	}
	response := u.Message(true, "Ok")
	response["race"] = race
	return response
}

func (entry *Race) IsClosed() bool {
	var race Race
	_, err := DBM.Query(&race, `SELECT * from races where status=? and id=?`, "Closed", entry.ID)
	if err != nil {
		panic(err)
		return false
	}
	fmt.Println(race)
	if race.ID != 0 {
		return true
	}
	return false
}
func (entry *Race) KeyIn() map[string]interface{} {
	var item Race
	if entry.Result == "" {
		return u.Message(false, "Please enter result!")
	}
	if !entry.IsClosed() {
		return u.Message(false, "Make sure the race status is CLOSED")
	}
	_, errdb := DBM.Query(&item, `SELECT * from races where id=?`, entry.ID)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	item.Status = "Complete"
	item.Remarks = "Race Complete"
	item.Result = entry.Result
	item.Completed = time.Now()
	errdb = DBM.Update(&item)
	if errdb != nil {
		panic(errdb)
		return u.Message(false, "Failed to create account, connection error")
	}
	item.ParseWinner()
	//func generate winners here
	go GenerateWinReports(entry.ID)

	response := u.Message(true, "KeyIn Successfully!")
	response["race"] = item
	return response
}
func (entry *Race) UpdateLap() map[string]interface{} {
	var item Race
	_, errdb := DBM.Query(&item, `SELECT * from races where id=?`, entry.ID)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	if item.Status == "Complete" {
		return u.Message(false, "Cannot update a completed race!")
	}
	item.Lap = entry.Lap

	errdb = DBM.Update(&item)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}

	//func generate winners here

	response := u.Message(true, "Lap Updated")
	return response
}

func (entry *Race) Update() map[string]interface{} {
	var item Race
	_, errdb := DBM.Query(&item, `SELECT * from races where id=?`, entry.ID)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}
	if item.Status == "Complete" {
		return u.Message(false, "Cannot update a completed race!")
	}
	fmt.Println(item)
	item.Status = entry.Status
	item.Remarks = entry.Remarks
	if item.Status == "Closed" {
		//genchan := make(chan int, 1)

		fmt.Println("closed race")
		go GenerateSalesReports(entry.ID)

		item.Odds = item.GetOddsVmix()
	}
	errdb = DBM.Update(&item)
	if errdb != nil {
		return u.Message(false, errdb.Error())
	}

	//func generate winners here

	response := u.Message(true, "Race updated successfully!")
	response["race"] = &item
	return response
}
