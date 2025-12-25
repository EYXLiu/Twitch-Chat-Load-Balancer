// decoder
//	decodes a message (converts from raw string into event)

package stream

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

func DecodeMessage(raw string) (*Event, error) {
	irc := DecodeIRCMessage(raw)

	switch irc.Command {
	case "PRIVMSG":
		payload := ChatEvent{
			Channel:   strings.TrimPrefix(irc.Params[0], "#"),
			User:      strings.Split(irc.Prefix, "!")[0],
			Message:   irc.Params[1],
			Timestamp: time.Now(),
		}
		data, _ := json.Marshal(payload)
		return &Event{
			Type: EventChat,
			Data: data,
		}, nil
	case "USERNOTICE":
		msg := ""
		if len(irc.Params) > 1 {
			msg = irc.Params[1]
		}
		payload := UserNoticeEvent{
			Channel:    strings.TrimPrefix(irc.Params[0], "#"),
			User:       irc.Tags["login"],
			Message:    msg,
			NoticeType: irc.Tags["msg-id"],
			SystemMsg:  strings.ReplaceAll(irc.Tags["system-msg"], `\s`, " "),
			Timestamp:  time.Now(),
		}
		data, _ := json.Marshal(payload)
		return &Event{
			Type: EventUserNotice,
			Data: data,
		}, nil
	case "CLEARCHAT":
		data, _ := json.Marshal(irc)
		return &Event{
			Type: EventClearChat,
			Data: data,
		}, nil
	default:
		return nil, errors.New("unsupported IRC command")
	}
}

func DecodeIRCMessage(raw string) IRCMessage {
	msg := IRCMessage{
		Tags: make(map[string]string),
	}

	// tags
	if strings.HasPrefix(raw, "@") {
		parts := strings.SplitN(raw, " ", 2)
		for tag := range strings.SplitSeq(parts[0][1:], ";") {
			kv := strings.SplitN(tag, "=", 2)
			if len(kv) == 2 {
				msg.Tags[kv[0]] = kv[1]
			}
		}
		raw = parts[1]
	}

	// prefix
	if strings.HasPrefix(raw, ":") {
		parts := strings.SplitN(raw[1:], " ", 2)
		msg.Prefix = parts[0]
		raw = parts[1]
	}
	// command + parts
	parts := strings.Split(raw, " ")
	msg.Command = parts[0]

	for i := 1; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], ":") {
			msg.Params = append(msg.Params, strings.Join(parts[i:], " ")[1:])
		}
		msg.Params = append(msg.Params, parts[i])
	}

	return msg
}
