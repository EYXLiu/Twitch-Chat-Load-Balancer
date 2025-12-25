// worker pool
// 	worker class that takes an event and pushes it to the atomic classes
//	calls counter, window, and cache
//	Submit
//		pushes a message into the workerpool input channel for workers to read
// 	worker
//		handles the messages and pushes to cache/window/counter

package processor

import (
	"encoding/json"
	"log"
	"sync"
	"tc/internal/analytics"
	"tc/internal/cache"
	"tc/internal/metrics"
	"tc/internal/stream"
)

type WorkerPool struct {
	workers int
	input   chan *stream.Event
	counter *metrics.Counter
	window  *analytics.Window
	cache   *cache.RedisCache
	wg      sync.WaitGroup
}

func WorkerPool_Init(workers int, counter *metrics.Counter, window *analytics.Window, cache *cache.RedisCache) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		input:   make(chan *stream.Event, 1000),
		counter: counter,
		window:  window,
		cache:   cache,
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	log.Printf("Started %d workers\n", p.workers)
}

func (p *WorkerPool) Submit(event *stream.Event) {
	select {
	case p.input <- event:
		break
	default:
		<-p.input
		p.input <- event
		log.Println("Error: worker pool full")
	}
}

func (p *WorkerPool) worker(id int) {
	_ = id
	defer p.wg.Done()
	for event := range p.input {
		switch event.Type {
		case stream.EventChat:
			var chat stream.ChatEvent
			if err := json.Unmarshal(event.Data, &chat); err != nil {
				log.Printf("worker %d: bad chat event: %v", id, err)
				return
			}
			p.counter.Inc()
			p.window.Add(&chat)
			p.cache.PushMessage(&chat)
		case stream.EventSystem:
			var notif stream.UserNoticeEvent
			if err := json.Unmarshal(event.Data, &notif); err != nil {
				log.Printf("worker %d: bad system event: %v", id, err)
				return
			}
			_ = notif
		default:

		}
	}
}

func (p *WorkerPool) Stop() {
	close(p.input)
	p.wg.Wait()
}
