package bus

import (
	"context"
	"encoding/json"
	"log"
	"tc/internal/stream"

	"github.com/redis/go-redis/v9"
)

type Consumer struct {
	rdb      *redis.Client
	group    string
	consumer string
	ctx      context.Context
}

func Consumer_Init(rdb *redis.Client, group, consumer string) *Consumer {
	return &Consumer{
		rdb:      rdb,
		group:    group,
		consumer: consumer,
		ctx:      context.Background(),
	}
}

func (c *Consumer) Start(handler func(*stream.ChatEvent)) {
	for {
		streams, err := c.rdb.XReadGroup(c.ctx, &redis.XReadGroupArgs{
			Group:    c.group,
			Consumer: c.consumer,
			Streams:  []string{"twitch:events", ">"},
			Count:    10,
			Block:    0,
		}).Result()

		if err != nil {
			log.Println("ERROR: redis read error:", err)
			continue
		}

		for _, s := range streams {
			for _, msg := range s.Messages {
				raw := msg.Values["data"].(string)

				var event stream.ChatEvent
				if err := json.Unmarshal([]byte(raw), &event); err != nil {
					continue
				}

				handler(&event)
				c.rdb.XAck(c.ctx, "twitch:events", c.group, msg.ID)
			}
		}
	}
}
