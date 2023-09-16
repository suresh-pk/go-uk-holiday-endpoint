package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var logger *log.Logger

func init() {
	// open the log file
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		logger.Println("Error opening log file:", err)
	}
	logger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type Event struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"-"`
	Notes   string    `json:"notes"`
	Bunting bool      `json:"bunting"`
}

type Division struct {
	Division string  `json:"division"`
	Events   []Event `json:"events"`
}

type HolidayResponse struct {
	EnglandAndWales Division `json:"england-and-wales"`
	Scotland        Division `json:"scotland"`
	NorthernIreland Division `json:"northern-ireland"`
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Title   string `json:"title"`
		DateStr string `json:"date"`
		Notes   string `json:"notes"`
		Bunting bool   `json:"bunting"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		logger.Println("Error:", err)
		return err
	}

	date, err := time.Parse("2006-01-02", tmp.DateStr)
	if err != nil {
		logger.Println("Error:", err)
		return err
	}

	e.Title = tmp.Title
	e.Date = date
	e.Notes = tmp.Notes
	e.Bunting = tmp.Bunting

	return nil
}

func main() {
	url := "https://www.gov.uk/bank-holidays.json"
	resp, err := http.Get(url)
	if err != nil {
		logger.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var data HolidayResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logger.Println("Error:", err)
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/holidays/{year}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		year := vars["year"]
		holidays := filterByYear(data, year)
		json.NewEncoder(w).Encode(holidays)
	}).Methods("GET")

	r.HandleFunc("/england-and-wales", func(w http.ResponseWriter, r *http.Request) {
		englandAndWalesData := filterByRegionBunting(data, "england-and-wales", false)
		json.NewEncoder(w).Encode(englandAndWalesData)
	}).Methods("GET")

	http.Handle("/", r)

	// Get all Bank holidays for a specific year and return title and date
	r.HandleFunc("/holidays-title-date/{year}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		year := vars["year"]
		holidays := filterByYear(data, year)
		var result []struct {
			Title string    `json:"title"`
			Date  time.Time `json:"date"`
		}

		for _, events := range holidays {
			for _, event := range events {
				result = append(result, struct {
					Title string    `json:"title"`
					Date  time.Time `json:"date"`
				}{
					Title: event.Title,
					Date:  event.Date,
				})
			}
		}

		json.NewEncoder(w).Encode(result)
	}).Methods("GET")

	// Start the HTTP server on localhost
	logger.Println("Server listening on :8080")
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

// Get holidays for all regions
func filterByYear(data HolidayResponse, year string) map[string][]Event {
	result := make(map[string][]Event)
	result["england-and-wales"] = filterEventsByYear(data.EnglandAndWales.Events, year)
	result["scotland"] = filterEventsByYear(data.Scotland.Events, year)
	result["northern-ireland"] = filterEventsByYear(data.NorthernIreland.Events, year)
	return result
}

func filterEventsByYear(events []Event, year string) []Event {
	var filteredEvents []Event

	for _, event := range events {
		eventYear := event.Date.Format("2006")
		if eventYear == year {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents
}

// Get all data of the region england and wales
func filterByRegionBunting(data HolidayResponse, region string, bunting bool) []Event {
	events := data.EnglandAndWales.Events

	var filteredEvents []Event

	for _, event := range events {
		if event.Bunting == bunting {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents
}
