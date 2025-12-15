package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tc/internal/bus"
	"tc/internal/config"
	"tc/internal/metrics"
	"tc/internal/stream"
	"tc/internal/twitch"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	client, _ := twitch.Connect()
	client.Join(cfg.TwitchChannel)

	msgQueue := make(chan *stream.ChatEvent, 1000)

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis,
	})

	producer := bus.Producer_Init(rdb)

	counter := metrics.Counter_Init()

	go func() {
		for event := range msgQueue {
			producer.Publish(event)
		}
	}()

	go client.Listen(func(raw string) {
		event, err := stream.DecodeIRCMessage(raw)
		if err != nil {
			return
		}

		select {
		case msgQueue <- event:
		default:
			counter.Inc()
			log.Println("Warning, dropping messages")
		}
	})

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down...")
	close(msgQueue)
	log.Printf("Dropped %d messages\n", counter.Get())
	client.Close()
}
