package common

import (
	"time"
	"fmt"
)

type runFunc func(interface{}) (error)

type AsyncJob struct {
	RunFunc   runFunc
	Name      string
	Type      string
	Parameter interface{}
}

func NewAsyncJobWorker(duration time.Duration) *AsyncJobWorker {
	asyncJobWorker := &AsyncJobWorker{
		closed:     false,
		closedChan: make(chan bool, 1),
		duration:   duration,
	}

	return asyncJobWorker
}

type AsyncJobWorker struct {
	queue      Queue
	duration   time.Duration
	closed     bool
	closedChan chan bool
}

func (myself *AsyncJobWorker) Start() {
	go myself.start()
}

func (myself *AsyncJobWorker) start() {
	var aysncJob *AsyncJob

	defer func() {
		err := recover()
		if nil != err {
			if nil != aysncJob {
				fmt.Printf("panic in job(type: %s; name: %s). error: %s", aysncJob.Type, aysncJob.Name, err)
			} else {
				fmt.Printf("panic error: %s", err)
			}

			myself.start()
		}
	}()

	for ; false == myself.closed; {
		for ; 0 < myself.queue.Length(); {
			item := myself.queue.Dequeue()
			aysncJob = item.(*AsyncJob)
			if nil == aysncJob {
				continue
			}

			aysncJob.RunFunc(aysncJob.Parameter)
		}

		time.Sleep(myself.duration)
	}

	myself.closedChan <- true
}

func (myself *AsyncJobWorker) Add(asyncJob *AsyncJob) {
	if myself.closed {
		return
	}

	myself.queue.Enqueue(asyncJob)
}

func (myself *AsyncJobWorker) Length() (int) {
	return myself.queue.Length()
}

func (myself *AsyncJobWorker) CloseAndWait() {
	myself.closed = true

	<-myself.closedChan
}
