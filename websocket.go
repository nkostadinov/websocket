package main

import (
	"github.com/nkostadinov/websocket/server"
	"net/http"
)


// This example demonstrates a trivial echo server.
func main() {
	//wesocket server
	webserver := server.NewServer("/")
	go webserver.Listen()
	//rest server
	restserver := server.NewRestServer(webserver);
	go restserver.ListenRest()

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
