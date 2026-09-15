package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	PORT = "7070"
)

func main() {
	if err := InitS3(); err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/zera/get_upload_url", GetPresignedUploadURL).Methods("POST")
	handler := cors.Default().Handler(router)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})
	handler = c.Handler(handler)
	fmt.Println("HYPERBALL server run at port: " + PORT)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
