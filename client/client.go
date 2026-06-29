package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"gollux/account"
	"gollux/auth"
)

const (
	PORT = "8000"
)

func main() {

	router := mux.NewRouter()
	// Routes consist of a path and a handler function.
	router.HandleFunc("/account/register", account.Register).Methods("POST")
	router.HandleFunc("/account/verify", account.Verify).Methods("POST")
	router.HandleFunc("/account/login_by_password", account.LoginByPassword).Methods("POST")
	router.HandleFunc("/account/transactions", account.GetTransactionsHistory).Methods("POST")
	router.HandleFunc("/account/get", account.GetAccount).Methods("POST")
	router.HandleFunc("/account/resendcode", account.ResendCode).Methods("POST")
	router.HandleFunc("/account/forgot", account.ForgotCode).Methods("POST")
	router.HandleFunc("/account/balance", account.BalanceInquire).Methods("POST")
	router.HandleFunc("/account/qry", account.PlayerCustomQry).Methods("POST")
	router.HandleFunc("/credits/buy", account.DoBuyCredits).Methods("POST")
	router.HandleFunc("/credits/cashout", account.DoCashout).Methods("POST")
	router.HandleFunc("/credits/cashout_cancel", account.UpdateCashout).Methods("POST")
	router.HandleFunc("/credits/get_cashout", account.GetCashout).Methods("POST")
	router.HandleFunc("/credits/callback", account.Callback).Methods("POST")
	router.HandleFunc("/credits/maya_callback", account.MayaCallback).Methods("POST")

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
	fmt.Println("HYPERBALL server run at port: " + PORT)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+PORT, handler))

}
