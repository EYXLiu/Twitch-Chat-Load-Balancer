// bus producer
//	receives an event and publishes it to a redis stream

package bus

import (
	"context"
	"encoding/json"
	"tc/internal/stream"
	"tc/internal/ws"

	"github.com/redis/go-redis/v9"
)

type Producer struct {
	rdb *redis.Client
	ctx context.Context
	hub *ws.Hub
}

func Producer_Init(rdb *redis.Client, hub *ws.Hub) *Producer {
	return &Producer{
		rdb: rdb,
		ctx: context.Background(),
		hub: hub,
	}
}

func (p *Producer) Publish(event *stream.ChatEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := "[" + event.Channel + "]" + event.User + ": " + event.Message
	p.hub.Broadcast <- []byte(msg)

	return p.rdb.XAdd(p.ctx, &redis.XAddArgs{
		Stream: "twitch:events",
		MaxLen: 100_000,
		Approx: true,
		Values: map[string]interface{}{
			"data": data,
		},
	}).Err()
}
