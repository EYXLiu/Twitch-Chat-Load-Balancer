package main

import (
	"log"
	"tc/internal/config"
	"tc/internal/stream"
	"tc/internal/twitch"
)

func main() {
	cfg := config.Load()

	client, err := twitch.Connect()
	if err != nil {
		log.Fatal(err)
	}

	client.Join(cfg.TwitchChannel)

	log.Println("Connected to twitch chat")

	client.Listen(func(raw string) {
		event, err := stream.DecodeIRCMessage(raw)
		if err != nil {
			log.Println("error:", err)
			return
		}

		log.Println(event.Message)
	})
}
