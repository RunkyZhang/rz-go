package common

import (
	"sync"
)

type RunFunc func(interface{}) (error)

type AsyncJob struct {
	RunFunc   RunFunc
	Name      string
	Type      string
	Parameter interface{}
}

func NewAsyncJobWorker(workerCount int) (*AsyncJobWorker) {
	asyncJobWorker := &AsyncJobWorker{
		closed: false,
	}
	asyncJobWorker.workerCount = workerCount
	asyncJobWorker.channel = make(chan bool, asyncJobWorker.workerCount)

	return asyncJobWorker
}

type AsyncJobWorker struct {
	queue       Queue
	closed      bool
	started     bool
	channel     chan bool
	workerCount int
	lock        sync.Mutex
	semaphore   Semaphore
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
			if nil != currentAsyncJob {
				GetLogging().Error(value,
					"Panic on job(type: %s; name: %s) in goroutine(%d)",
					currentAsyncJob.Type,
					currentAsyncJob.Name,
					id)
			} else {
				GetLogging().Error(value, "Panic in goroutine(%d)", id)
			}

			myself.start(id)
		}
	}()

	for ; ; {
		for ; ; {
			item := myself.queue.Dequeue()
			var ok bool
			currentAsyncJob, ok = item.(*AsyncJob)
			if !ok {
				break
			}

			err := currentAsyncJob.RunFunc(currentAsyncJob.Parameter)
			if nil != err {
				GetLogging().Warn(err, "Failed to run job in goroutine(%d)", id)
			}
		}

		if myself.closed {
			break
		}

		myself.semaphore.Wait()
	}

	GetLogging().Info(nil, "The goroutine(%d) is closing", id)
	myself.channel <- true
}

func (myself *AsyncJobWorker) Add(asyncJob *AsyncJob) {
	if myself.closed {
		return
	}

	myself.queue.Enqueue(asyncJob)
	myself.semaphore.Release()
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
	myself.semaphore.Release()

	for i := 0; i < myself.workerCount; i++ {
		<-myself.channel
	}
	close(myself.channel)
}
