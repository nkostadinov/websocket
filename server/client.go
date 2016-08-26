package server

import (
	"golang.org/x/net/websocket"
	"log"
	"encoding/json"
	"io"
)

// Chat client.
type Client struct {
	ws       *websocket.Conn
	server   *Server
	ch       chan *Message
	events   chan *Event
	done     chan bool
	channels []string
	session	 string
}

// write channel buffer size
const channelBufSize = 1000

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server, session string) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	} else if server == nil {
		panic("server cannot be nil")
	}

	ch := make(chan *Message, channelBufSize)
	done := make(chan bool)
	channels := make([]string, 0)
	events := make(chan *Event)

	return &Client{ws, server, ch, events, done, channels, session}
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
	//log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-self.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(self.ws, msg)
		//receive event from client
		case event := <-self.events:
			log.Println("Event: ", event)
			if(event.Event == "join") {
				self.server.joinChannel(event.Data, self)
			}
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
	//log.Println("Listening read from client")
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
		default:
			var in []byte
			if err := websocket.Message.Receive(self.ws, &in); err != nil {
				if err == io.EOF {
					self.done <- true
					return
				} else if err != nil {
					panic(err)
				}
			}
			log.Println("Recieved: ", string(in))

			var event Event
			if err := json.Unmarshal(in, &event); err != nil {
				panic(err)
			}
			self.events <- &event

		}
	}
}
