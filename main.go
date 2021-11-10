package main

import (
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var db, _ = gorm.Open("mysql", "root:root@/flip?charset=utf8&parseTime=True&loc=Local")

type TodoItemModel struct {
	Id          int `gorm:"primary_key"`
	Description string
	Completed   bool
}

func status(w http.ResponseWriter, r *http.Request) {
	log.Info("API Status is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
} // status

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
} // init

func main() {
	// Close our database connection when our main() function is returned.
	defer db.Close()

	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})

	// This is our router... similar to Laravel and Django ain't it?
	router := mux.NewRouter()
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/crawl", flips).Methods("GET")
	router.HandleFunc("/gas", average_regular_gas_price).Methods("GET")

	// Wrapping CORS handler around our existing application.
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	// Start the Server and listen on port 80
	http.ListenAndServe(":8000", handler)
	log.Info("Flip Server Started!")
} // main
