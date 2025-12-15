package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tc/internal/bus"
	"tc/internal/config"
	"tc/internal/handlers"
	"tc/internal/metrics"
	"tc/internal/twitch"
	"tc/internal/worker"
	"tc/internal/ws"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	r := gin.Default()

	client, _ := twitch.Connect()
	client.Join(cfg.TwitchChannel)

	msgQueue := make(chan string, 1000)

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis,
	})

	hub := ws.Hub_Init()
	go hub.Run()

	producer := bus.Producer_Init(rdb, hub)
	counter := metrics.Counter_Init()

	var workers []*worker.Worker
	done := make(chan int)
	maxWorkers := 5
	queueIncrement := 0.2
	idleTimeout := 5 * time.Second

	w := &worker.Worker{Id: 0, Quit: make(chan struct{})}
	workers = append(workers, w)
	go worker.RunWorker(0, msgQueue, producer, idleTimeout, done)

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		nextWorkerID := len(workers)

		for {
			select {
			case <-ticker.C:
				load := float64(len(msgQueue)) / float64(cap(msgQueue))
				desiredWorkers := min(int(load/queueIncrement)+1, maxWorkers)

				for len(workers) < desiredWorkers {
					w := &worker.Worker{Id: nextWorkerID, Quit: make(chan struct{})}
					workers = append(workers, w)
					nextWorkerID++
					go worker.RunWorker(w.Id, msgQueue, producer, idleTimeout, done)
					log.Printf("Spawning worker %d", w.Id)
				}
			case workerId := <-done:
				for i, w := range workers {
					if w.Id == workerId {
						workers = append(workers[:i], workers[i+1:]...)
						log.Printf("Worker %d terminated", workerId)
						break
					}
				}
			}
		}
	}()

	go client.Listen(func(raw string) {
		select {
		case msgQueue <- raw:
		default:
			counter.Inc()
			log.Println("Warning, dropping messages")
		}
	})

	handlers.AdminStatsHandler(r, counter, &workers, msgQueue)
	handlers.HealthHandler(r)
	handlers.WebsocketHandler(r, hub)

	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Gin server failed: %v", err)
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
	close(msgQueue)
	log.Printf("Dropped %d messages\n", counter.Get())
	client.Close()
}
