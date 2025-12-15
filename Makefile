SHELL := /bin/bash

ifndef TWITCH_CHANNEL
$(error TWITCH_CHANNEL is not set)
endif

ifndef REDIS
$(error REDIS is not set)
endif

export TWITCH_CHANNEL
export REDIS

run:
	@echo "reading $$TWITCH_CHANNEL chat"
	@echo "Starting redis consumer group"
	go run ./cmd/init || true
	@echo "Starting ingestor"
	go run ./cmd/ingestor & 
	@echo "Starting worker"
	go run ./cmd/worker

clean:
	redis-cli del twitch:messages
	redis-cli xgroup destroy twitch:events workers
	redis-cli del twitch:events