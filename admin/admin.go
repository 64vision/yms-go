package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"gollux/account"
	"gollux/auth"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	PORT = "8100"
)

func main() {
	fmt.Println("Starting administration")
	router := mux.NewRouter()
	router.HandleFunc("/admin/login", account.AdminLogin).Methods("POST")

	//

	router.HandleFunc("/admin/qry", account.CustomQry).Methods("POST")

	router.HandleFunc("/admin/player_stats", account.PlayerStats).Methods("GET")
	router.HandleFunc("/admin/settlements", account.GetSettlements).Methods("POST")
	router.HandleFunc("/admin/acct_settlement", account.GetAccountSettlement).Methods("POST")
	router.HandleFunc("/admin/update_settlement", account.UpdateSettlement).Methods("POST")
	router.HandleFunc("/admin/update_cashout", account.UpdateCashout).Methods("POST")
	router.HandleFunc("/admin/player_location", account.GetPlayersLocation).Methods("GET")
	//router.HandleFunc("/send/message", tools.SendMessage).Methods("POST")
	//router.HandleFunc("/message/list", tools.ListMessages).Methods("GET")

	//Credits path
	router.HandleFunc("/credits/topup", account.AddRequest).Methods("POST")

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
	//go BuyingTicker()
	/*--------------------------------------------------
		Run Server
	-----------------------------------------------------*/
	fmt.Println("HYPERBALL server run at port: " + PORT)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
