package server

import (
	"log"
	"net/http"
	"io/ioutil"
	"io"
)

type Rest struct {
	sendAll      chan *Message
}

func NewRestServer(server *Server) *Rest {
	return &Rest{server.sendAll}
}

func  (self *Rest) PostOnly(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h(w, r)
			return
		}
		http.Error(w, "post only", http.StatusMethodNotAllowed)
	}
}

func (self *Rest) restHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	channel := r.URL.Query().Get("channel")
	msg := &Message{ Channel: channel, Body: string(body) };
	self.sendAll <- msg;
	log.Printf("body: %s, channel: %s", body, channel);
}

func (self *Rest) ListenRest() {

	log.Println("Listening server(REST)...")

	http.HandleFunc("/rest", self.PostOnly(self.restHandler))
}
