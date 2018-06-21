package common

import (
	"time"
	"fmt"
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
	var currentAysncJob *AsyncJob

	//defer func() {
	//	value := recover()
	//	if nil != value {
	//		if nil != myself.defaultAsyncJob {
	//			fmt.Printf(
	//				"panic on job(type: %s; name: %s) in goroutine(%d). error: %s\n",
	//				myself.defaultAsyncJob.Type,
	//				myself.defaultAsyncJob.Name,
	//				id,
	//				fmt.Sprintln(value))
	//		} else if nil != currentAysncJob {
	//			fmt.Printf("panic on job(type: %s; name: %s) in goroutine(%d). error: %s\n",
	//				currentAysncJob.Type,
	//				currentAysncJob.Name,
	//				id,
	//				fmt.Sprintln(value))
	//		} else {
	//			fmt.Printf("panic in goroutine(%d). error: %s\n", id, fmt.Sprintln(value))
	//		}
	//
	//		myself.start(id)
	//	}
	//}()

	for ; ; {
		if nil != myself.defaultAsyncJob {
			myself.defaultAsyncJob.RunFunc(myself.defaultAsyncJob.Parameter)
		} else {
			for ; 0 < myself.queue.Length(); {
				item := myself.queue.Dequeue()
				var ok bool
				currentAysncJob, ok = item.(*AsyncJob)
				if !ok {
					continue
				}

				err := currentAysncJob.RunFunc(currentAysncJob.Parameter)
				if nil != err {
					fmt.Printf("failed to run job in goroutine(%d). error: %s\n", id, err)
				}
			}
		}

		if myself.closed {
			break
		}

		time.Sleep(myself.duration)
	}

	fmt.Printf("the goroutine(%d) is closing\n", id)
	myself.closedChan <- true
	fmt.Printf("the goroutine(%d) is closed\n", id)
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
