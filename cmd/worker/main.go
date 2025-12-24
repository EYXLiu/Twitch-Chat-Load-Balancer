// worker main
// 	starts the main worker goroutine
// 	sets up counter and window tracker
//	attaches to the pubsub bus
// 	prints every 5 seconds

package main

import (
	"fmt"
	"log"
	"os"
	"tc/internal/analytics"
	"tc/internal/bus"
	"tc/internal/cache"
	"tc/internal/config"
	"tc/internal/metrics"
	"tc/internal/processor"
	"tc/internal/stream"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis,
	})

	counter := metrics.Counter_Init()
	window := analytics.Window_Init(5 * time.Second)
	cache := cache.RedisCache_Init(rdb)

	pool := processor.WorkerPool_Init(4, counter, window, cache)
	pool.Start()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "worker-unknown"
	}
	consumerName := fmt.Sprintf("%s-%d", hostname, os.Getpid())
	consumer := bus.Consumer_Init(rdb, "workers", consumerName)

//	function just submits the event to the pool 
	go consumer.Start(func(event *stream.ChatEvent) {
		pool.Submit(event)
	})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			log.Println("Messages: ", counter.Get())
			log.Println("Messages last 5: ", window.Count())
		}
	}()

	select {}
}
