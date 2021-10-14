// client is used to handling client actions

package websockets

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	//Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	//Time allowed to read the next pong message from the peer. pong ensures that the peer is still connected
	pongWait = 60 * time.Second

	//Send pings to peer with this period. must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	//Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//Client is a middle man between the websocket connection and the hub
type Client struct {
	hub *Hub

	//The websocket connection
	conn *websocket.Conn

	//Buffered channel of outband messages
	send chan []byte
}

//readPump pumps messages from the websocket connection to the hub.
//
//the application runs readPump in a per-connection goroutine. The applicaiton ensures that there is at most
//one reader on a connection by executing all reads from this goroutine
func (c *Client) readPump() {
	defer func() { //runs clean up after the app is finished
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)              //sets the maximum message size which a user can send through a websocket
	c.conn.SetReadDeadline(time.Now().Add(pongWait)) //sets the deadline for the read to finish
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1)) // trims the whitespaces around and replaces newline with space
		c.hub.broadcast <- message
	}
}

//writePump pums message from the hub to the websocket connection.
//
//A go routine running writePump is started for each connection.
//the application ensures taht there is at most one writer to a connection by executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//the hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			//Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//serveWs handles websocket requests from the peer.
func ServeWs(id string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	hub := NewHub(id) //creates a new hub struct from hub.go file

	go hub.Run() //calls run method defined in hub.go file

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	//Allow collection of memory referenced by the caller by doing all work in new goroutines
	go client.writePump()
	go client.readPump()
}
