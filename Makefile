SHELL := /bin/bash

ifeq ($(strip $(TWITCH_CHANNELS)),)
$(error TWITCH_CHANNELS is not set (expected: channel1,channel2,channel3))
endif

REDIS ?= localhost:6379

export TWITCH_CHANNELS
export REDIS

all: run

run:
	@echo "reading $$TWITCH_CHANNELS chat"
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

.PHONY: all run clean