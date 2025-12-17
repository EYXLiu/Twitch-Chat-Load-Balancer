package twitch

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func Connect() (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(
		"wss://irc-ws.chat.twitch.tv:443",
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) send(msg string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (c *Client) Auth() error {
	c.send("NICK justinfan12345")
	return nil
}

func (c *Client) Join(channel string) error {
	c.send("JOIN #" + channel)
	return nil
}

func (c *Client) Send(message string) error {
	c.send(message)
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
