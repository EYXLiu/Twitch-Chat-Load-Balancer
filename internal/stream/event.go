package stream

import "time"

type EventType string

const (
	EventChat       EventType = "chat"
	EventUserNotice EventType = "usernotice"
	EventClearChat  EventType = "clearchat"
	EventSystem     EventType = "system"
)

type Event struct {
	Type EventType
	Data any
}

type IRCMessage struct {
	Tags    map[string]string
	Prefix  string
	Command string
	Params  []string
}

type ChatEvent struct {
	Channel   string
	User      string
	Message   string
	Timestamp time.Time
}
