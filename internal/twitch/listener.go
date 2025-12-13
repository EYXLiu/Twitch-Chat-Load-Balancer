package twitch

import "log"

func (c *Client) Listen(onMessage func(string)) {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("twitch read error:", err)
			return
		}
		onMessage(string(msg))
	}
}
