package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tc/internal/analytics"
	"tc/internal/cache"
	"tc/internal/config"
	"tc/internal/metrics"
	"tc/internal/processor"
	"tc/internal/stream"
	"tc/internal/twitch"
	"time"
)

func main() {
	cfg := config.Load()

	rawCh := make(chan string, 1000)

	client, _ := twitch.Connect()
	client.Join(cfg.TwitchChannel)

	counter := metrics.Counter_Init()
	window := analytics.Window_Init(5 * time.Second)
	cache := cache.RedisCache_Init(cfg.Redis)
	pool := processor.WorkerPool_Init(4, counter, window, cache)
	pool.Start()

	go client.Listen(func(raw string) {
		rawCh <- raw
	})

	go func() {
		for raw := range rawCh {
			event, err := stream.DecodeIRCMessage(raw)
			if err != nil {
				continue
			}
			pool.Submit(event)
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			log.Println("Messages: ", counter.Get())
			log.Print("Messages last 5: ", window.Count())
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down...")

	client.Close()
	close(rawCh)
	pool.Stop()
}
