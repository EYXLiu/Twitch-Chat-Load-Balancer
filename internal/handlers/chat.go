// chat handler
// 	chat ws endpoint
// 	streams all chat messages received

package handlers

import (
	"net/http"
	"tc/internal/twitch"
	"tc/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebsocketHandler(router *gin.Engine, hub *ws.Hub, client *twitch.Client) {
	router.GET("/chat", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		hub.AddClient(conn, client)
	})
}
