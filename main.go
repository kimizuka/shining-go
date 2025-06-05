package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/akualab/dmx"
)

var timer *time.Timer

func shining() {
	dmx, err := dmx.NewDMXConnection("/dev/tty.usbserial-EN437503")
	if err != nil {
		log.Fatal(err)
	}

	dmx.SetChannel(2, 100)
	dmx.Render()

	if timer != nil {
		timer.Stop()
	}

	timer = time.AfterFunc(1*time.Second, func() {
		dmx.SetChannel(2, 0)
		dmx.Render()
		dmx.Close()
	})
}

func shiningGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		response := struct {
			Message string `json:"message"`
		}{
			Message: "GET",
		}

		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		response := struct {
			Message string `json:"message"`
		}{
			Message: "POST",
		}

		shining()
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./public"))

	http.Handle("/", fs)
	http.HandleFunc("/api/shining", shiningGetHandler)

	port := ":8080"
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
