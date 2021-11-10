package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type AverageRegularGasPrice struct {
	Average string `json:"average"`
} // AverageRegularGasPrice

func average_regular_gas_price(w http.ResponseWriter, r *http.Request) {
	// Instantiate default collector
	freeCycleCollector := colly.NewCollector()

	freeCycleCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Requested: ", r.URL.String())
	})

	freeCycleCollector.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Request.URL, "Responded!")
	})

	// Parse throught the HTML to get all posts
	freeCycleCollector.OnHTML("p.numb", func(e *colly.HTMLElement) {
		// Parse out National Gas Price Average
		var ave AverageRegularGasPrice
		ave.Average = strings.Join(strings.Fields(strings.TrimSpace(e.DOM.Text())), " ")

		// Respond with Post JSON
		log.Info("API's Gas Price has a response!")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ave)
	})

	freeCycleCollector.Visit("https://gasprices.aaa.com/")
} // average_regular_gas_price

// func main() {
// 	// Instantiate default collector
// 	freeCycleCollector := colly.NewCollector()

// 	freeCycleCollector.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Requested: ", r.URL.String())
// 	})

// 	freeCycleCollector.OnResponse(func(r *colly.Response) {
// 		fmt.Println(r.Request.URL, "Responded!")
// 	})

// 	// Parse throught the HTML to get all posts
// 	freeCycleCollector.OnHTML("table.table-mob", func(e *colly.HTMLElement) {

// 		// Parse out the headers
// 		gas_types := make([]string, 0, 5)
// 		html_headers := e.DOM.ChildrenFiltered("thead").ChildrenFiltered("tr").ChildrenFiltered("th").Nodes
// 		for _, element := range html_headers {
// 			if element != nil && element.FirstChild != nil {
// 				fmt.Printf("%+v\n", element.FirstChild.Data)
// 				gas_types = append(gas_types, element.FirstChild.Data)
// 			} // if
// 		} // for
// 		fmt.Printf("%v", gas_types)

// 		// Retrieve the JSON stored in the fc-data html element's ":data" attribute
// 		// attrVal, _ := e.DOM.Attr(":data")

// 		// Stored the JSON data into a struct
// 		// var freecycle FreeCyleJSON
// 		// json.Unmarshal([]byte(attrVal), &freecycle)

// 		// Store db.Create(&freecycle.Posts)

// 		// Respond with Post JSON
// 		// log.Info("API's Crawler has a response!")
// 		// w.Header().Set("Content-Type", "application/json")
// 		// json.NewEncoder(w).Encode(freecycle.Posts)
// 	})

// 	freeCycleCollector.Visit("https://gasprices.aaa.com/")
// } // average_gas_prices
