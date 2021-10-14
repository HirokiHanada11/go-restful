// hub is used to manage the active clients and broadcasting of the messages

package websockets

import "log"

type Hub struct {
	//map of all the registered clients
	clients map[*Client]bool //Client struct is defined in client.go

	//inbound messages from the clients
	broadcast chan []byte

	//register requests from the clients
	register chan *Client

	//unregister requests from the clients
	unregister chan *Client
}

var Hubs = make(map[string]Hub) //map of all the hubs

/*
* takes hub id retrieved from the request parameter as an argument
* returns the pointer to the newly created hub or exisiting hub
 */
func NewHub(id string) *Hub {
	if hub, ok := Hubs[id]; ok {
		return &hub
	} else {
		log.Println("Creating new hub")
		hub := Hub{
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
		}
		Hubs[id] = hub
		return &hub
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: //if there is a data coming through the register channel
			log.Println("New Client registered")
			h.clients[client] = true //activates the client by setting the client key in clients map to true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client) //removes entry from a map with a key that matches the second arg
				close(client.send)
			}
		case message := <-h.broadcast:
			//if the hub struct includes message in broadcast channel, it sends the message to all the clients
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
