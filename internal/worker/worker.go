// worker
// 	RunWorker
//		auto stops worker when idle for idleTimeout
// 		publishes the decoded message to the producer bus

package worker

import (
	"log"
	"tc/internal/bus"
	"tc/internal/stream"
	"time"
)

type Worker struct {
	Id   int
	Quit chan struct{}
}

func RunWorker(id int, msgQueue chan string, producer *bus.Producer, idleTimeout time.Duration, done chan<- int) {
	idleTimer := time.NewTimer(idleTimeout)
	defer idleTimer.Stop()

	for {
		select {
		case raw, ok := <-msgQueue:
			if !ok {
				return
			}
			event, err := stream.DecodeMessage(raw)
			if err != nil {
				continue
			}
			if err := producer.Publish(event); err != nil {
				log.Printf("Worker %d: failed to publish %v", id, err)
				continue
			}
			idleTimer.Reset(idleTimeout)
		case <-idleTimer.C:
			done <- id
			return
		}
	}
}
