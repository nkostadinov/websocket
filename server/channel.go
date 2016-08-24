package server

import "log"


// Chat server.
type Channel struct {
	name         string
	clients      []*Client
	addClient    chan *Client
	removeClient chan *Client
	sendAll      chan *Message
}


func NewChannel(channel string) *Channel {
	clients := make([]*Client, 0)
	addClient := make(chan *Client)
	removeClient := make(chan *Client)
	sendAll := make(chan *Message)
	return &Channel{channel, clients, addClient, removeClient, sendAll}
}

func (self *Channel) Listen() {
	log.Println("Listening ", self.name)

	for {
		select {

		// Add new a client
		case c := <-self.addClient:
			log.Println("Added new client to channel")
			for _, cli := range self.clients {
				if cli == c {
					return
				}
			}
			self.clients = append(self.clients, c)
			log.Println("Now", len(self.clients), "clients connected in channel ", self.name)

		// remove a client
		case c := <-self.removeClient:
			log.Println("Remove client from channel", self.name)
			for i := range self.clients {
				if self.clients[i] == c {
					self.clients = append(self.clients[:i], self.clients[i + 1:]...)
					break
				}
			}

		// broadcast message for all clients
		case msg := <-self.sendAll:
			log.Println("Send all in channel:", msg)
			for _, c := range self.clients {
				c.Write() <- msg
			}
		}
	}

}

