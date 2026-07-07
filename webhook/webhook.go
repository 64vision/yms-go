package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"gollux/auth"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	PORT = "9700"
)

func main() {

	router := mux.NewRouter()
	// Routes consist of a path and a handler function.
	router.HandleFunc("/api/webhook", ApiWebhook).Methods("POST")

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
	fmt.Println("HYPZera webhook server run at port: " + PORT)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+PORT, handler))

}
