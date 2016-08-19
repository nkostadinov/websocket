package main

import (
	"./server"
	"net/http"
)

// This example demonstrates a trivial echo server.
func main() {
	//rest server
	go server.ListenRest()
	//wesocket server
	webserver := server.NewServer("/")
	go webserver.Listen()

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
