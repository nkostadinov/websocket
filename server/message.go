package server

import (
	"encoding/json"
)

type Message struct {
	Channel string `json:"channel"`
	Body    string `json:"body"`
}

func (self *Message) String() string {
	b, err := json.Marshal(self)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type Event struct {
	Event	string	`json:"_event"`
	Data	string	`json:"_data"`
}

func (self *Event) String() string {
	b, err := json.Marshal(self)
	if err != nil {
		panic(err)
	}
	return string(b)
}