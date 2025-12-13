package config

import (
	"os"
)

type Config struct {
	TwitchChannel string
	Redis         string
}

func Load() Config {
	return Config{
		TwitchChannel: os.Getenv("TWITCH_CHANNEL"),
		Redis:         os.Getenv("REDIS"),
	}
}
