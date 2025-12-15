package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tc/internal/bus"
	"tc/internal/config"
	"tc/internal/stream"
	"tc/internal/twitch"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	client, _ := twitch.Connect()
	client.Join(cfg.TwitchChannel)

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis,
	})

	producer := bus.Producer_Init(rdb)

	go client.Listen(func(raw string) {
		event, err := stream.DecodeIRCMessage(raw)
		if err != nil {
			return
		}
		producer.Publish(event)
	})

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down...")
	client.Close()
}
