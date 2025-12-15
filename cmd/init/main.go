package main

import (
	"tc/internal/bus"
	"tc/internal/config"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis,
	})

	bus.Redis_Init(rdb)
}
