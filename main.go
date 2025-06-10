package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akualab/dmx"
)

func onLight(level int) {
	dmx, err := dmx.NewDMXConnection("/dev/tty.usbserial-EN437503")
	if err != nil {
		log.Fatal(err)
	}

	dmx.SetChannel(2, byte(level))
	dmx.Render()
	dmx.Close()
}

func shiningGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var request struct {
			Level int `json:"level"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.Level < 0 || 100 < request.Level {
			http.Error(w, "Level must be between 0 and 100", http.StatusBadRequest)
			return
		}

		onLight(request.Level)

		response := struct {
			Level int `json:"level"`
		}{
			Level: request.Level,
		}

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
