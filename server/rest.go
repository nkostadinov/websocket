package server

import (
	"log"
	"net/http"
	"io/ioutil"
	"io"
)

func restHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	channel := r.URL.Query().Get("channel")
	log.Printf("body: %s, channel: %s", body, channel);
}

func ListenRest() {

	log.Println("Listening server(REST)...")

	http.HandleFunc("/rest", restHandler)
}
