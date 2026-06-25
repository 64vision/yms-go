package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"hyperball.com/account"
	"hyperball.com/auth"
)

const (
	PORT = "8200"
)

func main() {
	fmt.Println("Starting agent endpoint...")
	router := mux.NewRouter()
	router.HandleFunc("/account/login_by_password", account.LoginByPassword).Methods("POST")
	router.HandleFunc("/agent/qry", account.CustomQry).Methods("POST")
	router.HandleFunc("/agent/data_qry", Query).Methods("POST")
	router.HandleFunc("/account/qry", account.PlayerCustomQry).Methods("POST")
	router.HandleFunc("/account/balance", account.BalanceInquire).Methods("POST")
	router.HandleFunc("/credits/buy", account.DoBuyCredits).Methods("POST")
	router.HandleFunc("/credits/send", account.AddRequest).Methods("POST")
	router.HandleFunc("/credits/proceed", account.ProceedRequest).Methods("POST")
	router.HandleFunc("/get_access", GetAccessAPI).Methods("GET")
	router.HandleFunc("/agent/performance", GetPerformance).Methods("POST")

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

	/*--------------------------------------------------
		Run Server
	-----------------------------------------------------*/
	fmt.Println("HYPERBALL server agent run at port: " + PORT)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+PORT, handler))

}
