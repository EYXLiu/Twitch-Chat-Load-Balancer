package bus

import (
	"context"
	"encoding/json"
	"log"
	"tc/internal/stream"

	"github.com/redis/go-redis/v9"
)

type Producer struct {
	rdb *redis.Client
	ctx context.Context
}

func Producer_Init(rdb *redis.Client) *Producer {
	return &Producer{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (p *Producer) Publish(event *stream.ChatEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	log.Println(event.User, ":", event.Message)

	return p.rdb.XAdd(p.ctx, &redis.XAddArgs{
		Stream: "twitch:events",
		MaxLen: 100_000,
		Approx: true,
		Values: map[string]interface{}{
			"data": data,
		},
	}).Err()
}
