package stream

import "time"

type ChatEvent struct {
	Channel   string
	User      string
	Message   string
	Timestamp time.Time
}
