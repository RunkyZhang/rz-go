package common

import (
	"time"
	"sync"
)

type runFunc func(interface{}) (error)

type AsyncJob struct {
	RunFunc   runFunc
	Name      string
	Type      string
	Parameter interface{}
}

func NewAsyncJobWorker(workerCount int, duration time.Duration, defaultAsyncJob ...*AsyncJob) (*AsyncJobWorker) {
	asyncJobWorker := &AsyncJobWorker{
		closed:   false,
		duration: duration,
	}
	if nil == defaultAsyncJob || 0 == len(defaultAsyncJob) {
		asyncJobWorker.defaultAsyncJob = nil
	} else {
		asyncJobWorker.defaultAsyncJob = defaultAsyncJob[0]
	}
	asyncJobWorker.workerCount = workerCount
	asyncJobWorker.closedChan = make(chan bool, asyncJobWorker.workerCount)

	return asyncJobWorker
}

type AsyncJobWorker struct {
	queue           Queue
	duration        time.Duration
	closed          bool
	started         bool
	closedChan      chan bool
	workerCount     int
	defaultAsyncJob *AsyncJob
	lock            sync.Mutex
}

func (myself *AsyncJobWorker) Start() {
	if myself.started {
		return
	}

	myself.lock.Lock()
	defer myself.lock.Unlock()

	if myself.started {
		return
	}

	for i := 0; i < myself.workerCount; i++ {
		go myself.start(i)
	}

	myself.started = true
}

func (myself *AsyncJobWorker) start(id int) {
	var currentAsyncJob *AsyncJob

	defer func() {
		value := recover()
		if nil != value {
			if nil != myself.defaultAsyncJob {
				GetLogging().Error(value,
					"panic on job(type: %s; name: %s) in goroutine(%d)",
					myself.defaultAsyncJob.Type,
					myself.defaultAsyncJob.Name,
					id)
			} else if nil != currentAsyncJob {
				GetLogging().Error(value,
					"panic on job(type: %s; name: %s) in goroutine(%d)",
					currentAsyncJob.Type,
					currentAsyncJob.Name,
					id)
			} else {
				GetLogging().Error(value, "panic in goroutine(%d)", id)
			}

			myself.start(id)
		}
	}()

	for ; ; {
		if nil != myself.defaultAsyncJob {
			myself.defaultAsyncJob.RunFunc(myself.defaultAsyncJob.Parameter)
		} else {
			for ; 0 < myself.queue.Length(); {
				item := myself.queue.Dequeue()
				var ok bool
				currentAsyncJob, ok = item.(*AsyncJob)
				if !ok {
					continue
				}

				err := currentAsyncJob.RunFunc(currentAsyncJob.Parameter)
				if nil != err {
					GetLogging().Warn(err, "failed to run job in goroutine(%d)", id)
				}
			}
		}

		if myself.closed {
			break
		}

		time.Sleep(myself.duration)
	}

	GetLogging().Info(nil, "the goroutine(%d) is closing", id)
	myself.closedChan <- true
	GetLogging().Info(nil, "the goroutine(%d) is closed", id)
}

func (myself *AsyncJobWorker) Add(asyncJob *AsyncJob) {
	if myself.closed || nil != myself.defaultAsyncJob {
		return
	}

	myself.queue.Enqueue(asyncJob)
}

func (myself *AsyncJobWorker) QueueLength() (int) {
	return myself.queue.Length()
}

func (myself *AsyncJobWorker) WorkerCount() (int) {
	return myself.workerCount
}

func (myself *AsyncJobWorker) CloseAndWait() {
	if !myself.started {
		return
	}
	if myself.closed {
		return
	}

	myself.lock.Lock()
	defer myself.lock.Unlock()

	if !myself.started {
		return
	}
	if myself.closed {
		return
	}

	myself.closed = true

	for i := 0; i < myself.workerCount; i++ {
		<-myself.closedChan
	}
}
