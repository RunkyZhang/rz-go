package common

import (
	"time"
)

func NewAsyncJobTrigger(workerCount int, interval time.Duration, asyncJob AsyncJob) (*AsyncJobTrigger) {
	asyncJobTrigger := &AsyncJobTrigger{
		closed:   false,
		interval: interval,
		asyncJob: asyncJob,
	}
	asyncJobTrigger.ticker = time.NewTicker(asyncJobTrigger.interval)
	asyncJobTrigger.asyncJobWorker = NewAsyncJobWorker(workerCount)

	return asyncJobTrigger
}

type AsyncJobTrigger struct {
	closed         bool
	interval       time.Duration
	asyncJob       AsyncJob
	ticker         *time.Ticker
	asyncJobWorker *AsyncJobWorker
}

func (myself *AsyncJobTrigger) Start() {
	myself.asyncJobWorker.Start()

	go myself.start()
}

func (myself *AsyncJobTrigger) start() {
	for range myself.ticker.C {
		if myself.closed {
			break
		}

		count := myself.asyncJobWorker.WorkerCount()
		for i := 0; i < count; i++ {
			myself.asyncJobWorker.Add(&myself.asyncJob)
		}
	}
}

func (myself *AsyncJobTrigger) CloseAndWait() {
	myself.closed = true

	myself.asyncJobWorker.CloseAndWait()
}
