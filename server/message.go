package server

import (
	"encoding/json"
)

type Message struct {
	Channel string `json:"channel"`
	Body   string `json:"body"`
}

func (self *Message) String() string {
	b, err := json.Marshal(self)
	if err != nil {
		panic(err)
	}
	return string(b)
}
