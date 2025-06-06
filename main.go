package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akualab/dmx"
)

var isOn bool

func offLight() {
	dmx, err := dmx.NewDMXConnection("/dev/tty.usbserial-EN437503")
	if err != nil {
		log.Fatal(err)
	}

	dmx.SetChannel(2, 0)
	dmx.Render()
	dmx.Close()
}

func onLight() {
	dmx, err := dmx.NewDMXConnection("/dev/tty.usbserial-EN437503")
	if err != nil {
		log.Fatal(err)
	}

	dmx.SetChannel(2, 100)
	dmx.Render()
	dmx.Close()
}

func shiningGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		if !isOn {
			isOn = true
			onLight()
		} else {
			isOn = false
			offLight()
		}

		response := struct {
			IsOn bool `json:"isOn"`
		}{
			IsOn: isOn,
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
