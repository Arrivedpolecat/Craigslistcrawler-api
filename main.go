package main

import (
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

type TodoItemModel struct {
	Id          int `gorm:"primary_key"`
	Description string
	Completed   bool
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
} // init

func main() {
	// TODO: Run locally
	// port := "8000"

	// TODO: For production
	port := os.Getenv("PORT")

	// This is our router... similar to Laravel and Django ain't it?
	router := mux.NewRouter()
	router.HandleFunc("/crawl", freeItems).Methods("GET")
	router.HandleFunc("/gas", average_regular_gas_price).Methods("GET")

	// Wrapping CORS handler around our existing application.
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	// Start the Server and listen on port 80
	http.ListenAndServe(":"+port, handler)
	log.Info("Flip Server Started!")
} // main
