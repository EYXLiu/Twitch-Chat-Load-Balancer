// counter metric
// 	atomic counter to increment and get the total

package metrics

import "sync/atomic"

type Counter struct {
	total int64
}

func Counter_Init() *Counter {
	return &Counter{}
}

func (c *Counter) Inc() {
	atomic.AddInt64(&c.total, 1)
}

func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.total)
}
