// websocket hub for twitch
//	sets up the websocket and broadcasts
//	calls connection write message to twitch (eg. join, renick, send chat message)
//	Hub_Run
//		sends the messages in the Broadcast Hub channel to each Send Client channel
// 	Client_Run
//		sends the message from the Send Client channel to the websocket connection
//	Hub_AddClient
//		adds clients with the requested streamtypes as a Client struct
//	Hub_ReadLoop
//		reads messages from the websocket and sends them to twitch

package ws

import (
	"encoding/json"
	"sync"
	"tc/internal/stream"
	"tc/internal/twitch"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	twitch   *twitch.Client
	interest map[stream.EventType]bool
}

type Hub struct {
	clients   map[*Client]bool
	Broadcast chan *stream.Event
	lock      sync.Mutex
}

func Hub_Init() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		Broadcast: make(chan *stream.Event),
	}
}

func (h *Hub) Run() {
	for msg := range h.Broadcast {
		h.lock.Lock()
		for c := range h.clients {
			if !c.interest[msg.Type] {
				continue
			}

			data, _ := json.Marshal(msg)
			select {
			case c.send <- data:
			default:
				close(c.send)
				delete(h.clients, c)
			}
		}
		h.lock.Unlock()
	}
}

func (c *Client) Run(h *Hub) {
	defer c.conn.Close()

	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (h *Hub) AddClient(conn *websocket.Conn, twitch *twitch.Client, events []stream.EventType) {
	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 256),
		twitch:   twitch,
		interest: make(map[stream.EventType]bool),
	}

	for _, e := range events {
		client.interest[e] = true
	}

	h.lock.Lock()
	h.clients[client] = true
	h.lock.Unlock()

	go client.Run(h)
	go client.ReadLoop(h)
}

func (c *Client) ReadLoop(h *Hub) {
	defer func() {
		h.lock.Lock()
		delete(h.clients, c)
		h.lock.Unlock()
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		if err := c.twitch.Send(string(msg)); err != nil {
			c.send <- []byte("[TWITCH]error: " + err.Error())
		}
	}
}
