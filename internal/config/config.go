// config
// 	read config from environment
//	load into Config struct

package config

import (
	"os"
	"strings"
)

type Config struct {
	TwitchChannels []string
	Redis          string
}

func Load() Config {
	tcs := strings.Split(os.Getenv("TWITCH_CHANNELS"), ",")

	for i := range tcs {
		tcs[i] = strings.TrimSpace(tcs[i])
	}

	return Config{
		TwitchChannels: tcs,
		Redis:          os.Getenv("REDIS"),
	}
}
