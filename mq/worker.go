package mq

import (
	"sync"
)

func NewWorker(process func() error, numWorkers int) *Worker {
	return &Worker{
		numWorkers: numWorkers,
		process:    process,
	}
}

type Worker struct {
	numWorkers  int
	process     func() error
	wg          sync.WaitGroup
	stopChannel chan struct{}
}

func (mq *Worker) Stop() {
	for i := 0; i < mq.numWorkers; i++ {
		mq.stopChannel <- struct{}{}
	}
	mq.wg.Wait()
}

func (mq *Worker) Start() {
	mq.stopChannel = make(chan struct{}, 1)
	mq.wg.Add(mq.numWorkers)

	for i := 0; i < mq.numWorkers; i++ {
		go func(workerID int) {
			defer mq.wg.Done()

			for {
				select {
				case <-mq.stopChannel:
					return
				default:
					if err := mq.process(); err != nil {
						continue
					}
				}
			}
		}(i)
	}
}
