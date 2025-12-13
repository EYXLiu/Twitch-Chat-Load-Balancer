package cache

import (
	"context"
	"tc/internal/stream"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func RedisCache_Init(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisCache) PushMessage(event *stream.ChatEvent) error {
	key := "twitch:messages"
	value := event.User + ": " + event.Message
	return r.client.RPush(r.ctx, key, value).Err()
}
