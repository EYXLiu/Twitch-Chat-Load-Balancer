package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	lock      sync.Mutex
}

func Hub_Init() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		msg := <-h.Broadcast
		h.lock.Lock()
		for client := range h.clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
		h.lock.Unlock()
	}
}

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.lock.Lock()
	h.clients[conn] = true
	h.lock.Unlock()
}
