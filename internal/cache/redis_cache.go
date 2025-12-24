// redis cache functions
//	initialize redis
//	push a twitch message into twitch:messages

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

func RedisCache_Init(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisCache) PushMessage(event *stream.ChatEvent) error {
	key := "twitch:messages"
	value := "[" + event.Channel + "]" + event.User + ": " + event.Message
	return r.client.RPush(r.ctx, key, value).Err()
}
