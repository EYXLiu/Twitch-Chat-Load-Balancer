package processor

import (
	"log"
	"sync"
	"tc/internal/analytics"
	"tc/internal/cache"
	"tc/internal/metrics"
	"tc/internal/stream"
)

type WorkerPool struct {
	workers int
	input   chan *stream.ChatEvent
	counter *metrics.Counter
	window  *analytics.Window
	cache   *cache.RedisCache
	wg      sync.WaitGroup
}

func WorkerPool_Init(workers int, counter *metrics.Counter, window *analytics.Window, cache *cache.RedisCache) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		input:   make(chan *stream.ChatEvent, 1000),
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

func (p *WorkerPool) Submit(event *stream.ChatEvent) {
	p.input <- event
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()
	for event := range p.input {
		p.counter.Inc()
		p.window.Add(event)
		p.cache.PushMessage(event)
		_ = id
		_ = event
	}
}

func (p *WorkerPool) Stop() {
	close(p.input)
	p.wg.Wait()
}
