package server

import (
	"golang.org/x/net/websocket"
	"log"
)

// Chat client.
type Client struct {
	ws       *websocket.Conn
	server   *Server
	ch       chan *Message
	done     chan bool
	channels []string
}

// write channel buffer size
const channelBufSize = 1000

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	} else if server == nil {
		panic("server cannot be nil")
	}

	ch := make(chan *Message, channelBufSize)
	done := make(chan bool)
	channels := make([]string, 0)

	return &Client{ws, server, ch, done, channels}
}

// Get websocket connection.
func (self *Client) Conn() *websocket.Conn {
	return self.ws
}

// Get Write channel
func (self *Client) Write() chan <- *Message {
	return (chan <- *Message)(self.ch)
}

// Get done channel.
func (self *Client) Done() chan <- bool {
	return (chan <- bool)(self.done)
}

// Listen Write and Read request via chanel
func (self *Client) Listen() {
	go self.listenWrite()
	self.listenRead()
}

// Listen write request via chanel
func (self *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-self.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(self.ws, msg)

		// receive done request
		case <-self.done:
			self.server.RemoveClient() <- self
			for _, ch := range self.channels {
				log.Println("removing from ", ch)
				self.server.channels[ch].removeClient <- self
			}
			self.done <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (self *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-self.done:
			self.server.RemoveClient() <- self
			for _, ch := range self.channels {
				log.Println("removing from ", ch)
				self.server.channels[ch].removeClient <- self
			}
			self.done <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(self.ws, &msg)
			if err != nil {
				if(err.Error() == "EOF") {
					self.done <- true
				}

				log.Println("String not json.", err)
			}
		//} else {
		//	self.server.SendAll() <- &msg
		//}
		}
	}
}
