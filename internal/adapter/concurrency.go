package adapter

import (
  "context"
  "fmt"
  "log"
  "time"
)

// --- In-Memory Queue ---
type ChanQueue struct {
  ch chan interface{}
}

func NewChanQueue() *ChanQueue {
  // queue 100 slots long
  q := &ChanQueue{ch: make(chan interface{}, 100)}
  // worker to consume messages
  go func() {
    for msg := range q.ch {
      log.Printf("[Async Worker] Processed event: %+v\n", msg)
    }
  }()
  return q
}

func (q *ChanQueue) Publish(ctx context.Context, topic string, payload interface{}) error {
  // Wrapping payload to include topic
  msg := map[string]interface{}{
    "topic":     topic,
    "payload":   payload,
    "timestamp": time.Now(),
  }
  select {
  case q.ch <- msg:
    return nil
  default:
    return fmt.Errorf("queue full")
  }
}
