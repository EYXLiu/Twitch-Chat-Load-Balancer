package config

import (
	"os"
)

type Config struct {
	TwitchChannel string
}

func Load() Config {
	return Config{
		TwitchChannel: os.Getenv("TWITCH_CHANNEL"),
	}
}
