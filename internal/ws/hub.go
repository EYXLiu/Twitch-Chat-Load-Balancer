package ws

import (
	"sync"
	"tc/internal/twitch"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	twitch *twitch.Client
}

type Hub struct {
	clients   map[*Client]bool
	Broadcast chan []byte
	lock      sync.Mutex
}

func Hub_Init() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		Broadcast: make(chan []byte),
	}
}

func (h *Hub) Run() {
	for msg := range h.Broadcast {
		h.lock.Lock()
		for c := range h.clients {
			select {
			case c.send <- msg:
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

func (h *Hub) AddClient(conn *websocket.Conn, twitch *twitch.Client) {
	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		twitch: twitch,
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
