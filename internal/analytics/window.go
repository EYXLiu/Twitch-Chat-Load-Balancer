// analytics (5 second window)
//	creates a window for the last 5 seconds and returns analytics ()

package analytics

import (
	"sync"
	"tc/internal/stream"
	"time"
)

type Window struct {
	messages []*stream.ChatEvent
	mutex    sync.Mutex
	duration time.Duration
}

func Window_Init(duration time.Duration) *Window {
	return &Window{
		messages: make([]*stream.ChatEvent, 0),
		duration: duration,
	}
}

func (w *Window) Add(event *stream.ChatEvent) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.messages = append(w.messages, event)
	w.expireOld(event.Timestamp)
}

func (w *Window) expireOld(now time.Time) {
	cutoff := now.Add(-w.duration)
	i := 0
	for ; i < len(w.messages); i++ {
		if w.messages[i].Timestamp.After(cutoff) {
			break
		}
	}
	w.messages = w.messages[i:]
}

func (w *Window) GetMessages() []*stream.ChatEvent {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return append([]*stream.ChatEvent{}, w.messages...)
}

func (w *Window) Count() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return len(w.messages)
}
