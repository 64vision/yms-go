package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"zerasuite/bookings"

	"gollux/account"
	"gollux/auth"
	"gollux/dbconfig"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	PORT = "9900" //9900-msk
)

var DBM *pg.DB

func main() {
	fmt.Println("Starting administration")
	DBM = dbconfig.DBM
	router := mux.NewRouter()
	router.HandleFunc("/admin/login", account.AdminLogin).Methods("POST")

	//

	router.HandleFunc("/admin/qry", account.CustomQry).Methods("POST")
	router.HandleFunc("/exec/qry", account.ExecCustomQry).Methods("POST")
	// router.HandleFunc("/admin/player_stats", account.PlayerStats).Methods("GET")
	// router.HandleFunc("/admin/settlements", account.GetSettlements).Methods("POST")
	// router.HandleFunc("/admin/acct_settlement", account.GetAccountSettlement).Methods("POST")
	// router.HandleFunc("/admin/update_settlement", account.UpdateSettlement).Methods("POST")
	// router.HandleFunc("/admin/update_cashout", account.UpdateCashout).Methods("POST")
	// router.HandleFunc("/admin/player_location", account.GetPlayersLocation).Methods("GET")
	//router.HandleFunc("/send/message", tools.SendMessage).Methods("POST")
	//router.HandleFunc("/message/list", tools.ListMessages).Methods("GET")

	//Credits path
	router.HandleFunc("/credits/topup", account.AddRequest).Methods("POST")

	//Bookings path

	router.Use(auth.JwtAuthentication)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	handler := cors.Default().Handler(router)
	handler = c.Handler(handler)
	rand.Seed(time.Now().UnixNano())
	//sms.Send("09156033392", "test")
	go RunTicker()
	/*--------------------------------------------------
		Run Server
	-----------------------------------------------------*/
	fmt.Println("HYPERBALL server run at port: " + PORT)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}

func RunTicker() {
	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				NoShow()

			}
		}
	}()
}

// run no show
func NoShow() {
	fmt.Println("Running no show filter")
	var bookings []bookings.Booking
	res, err := DBM.Query(&bookings, `UPDATE bookings
SET status = 'Forfeited', remarks='No Show'
WHERE remarks IS NULL
  AND status = 'Active'
  AND (booking_date::date + slot_time::time)
      <= NOW() - INTERVAL '2 hour'
RETURNING *`)
	if err != nil {
		panic(err)
	}
	if res.RowsReturned() == 0 {
		return
	}
	for _, booking := range bookings {
		trans := &account.Transaction{}
		trans.Type = "refund"
		trans.AccountID = booking.ClientID
		trans.Amount = booking.DocsFee
		trans.CreatedAt = time.Now()
		trans.RefNo = booking.ID
		trans.Description = "Docs fee refund for forfeited no show booking. Ref ID: " + fmt.Sprint(booking.ID)
		trans.Add()
	}

}
