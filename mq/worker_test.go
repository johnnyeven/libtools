package mq

import (
	"sync/atomic"
	"testing"
)

func TestWorker(t *testing.T) {
	count := int64(0)

	process := func() error {
		c := atomic.LoadInt64(&count)
		atomic.StoreInt64(&count, c+1)
		return nil
	}

	worker := NewWorker(process, 2)
	go worker.Start()

	for {
		c := atomic.LoadInt64(&count)
		if c > 10 {
			worker.Stop()
			break
		}
	}
}
