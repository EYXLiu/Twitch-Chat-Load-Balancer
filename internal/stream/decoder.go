package stream

import (
	"errors"
	"strings"
	"time"
)

type ChatEvent struct {
	Channel   string
	User      string
	Message   string
	Timestamp time.Time
}

func DecodeIRCMessage(raw string) (*ChatEvent, error) {
	if !strings.Contains(raw, "PRIVMSG") {
		return nil, errors.New("not a chat message")
	}

	parts := strings.Split(raw, "!")
	if len(parts) < 2 {
		return nil, errors.New("invalid format")
	}
	user := parts[0][1:]

	msgParts := strings.SplitN(raw, ":", 3)
	if len(msgParts) < 3 {
		return nil, errors.New("no message found")
	}
	message := strings.TrimSpace(msgParts[2])

	channelParts := strings.Split(raw, "PRIVMSG #")
	if len(channelParts) < 2 {
		return nil, errors.New("no channel found")
	}
	channel := strings.Split(channelParts[1], " ")[0]

	return &ChatEvent{
		Channel:   channel,
		User:      user,
		Message:   message,
		Timestamp: time.Now(),
	}, nil
}
