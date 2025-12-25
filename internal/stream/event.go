// event
//	event types and structs 

package stream

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	EventChat       EventType = "chat"
	EventUserNotice EventType = "usernotice"
	EventClearChat  EventType = "clearchat"
	EventSystem     EventType = "system"
)

type Event struct {
	Type EventType
	Data json.RawMessage
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

type UserNoticeEvent struct {
	Channel    string
	User       string
	Message    string
	NoticeType string
	SystemMsg  string
	Timestamp  time.Time
}
