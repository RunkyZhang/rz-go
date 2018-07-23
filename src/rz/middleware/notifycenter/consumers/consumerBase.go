package consumers

import (
	"time"
	"runtime"
	"rz/middleware/notifycenter/common"
)

type consumerBase struct {
	name           string
	asyncJobWorker *common.AsyncJobWorker
	runFunc        common.RunFunc
}

func (myself *consumerBase) Start(interval time.Duration) (error) {
	asyncJob := &common.AsyncJob{
		Name:    myself.name,
		Type:    "Consumer",
		RunFunc: myself.runFunc,
	}
	myself.asyncJobWorker = common.NewAsyncJobWorker(runtime.NumCPU(), interval, asyncJob)
	myself.asyncJobWorker.Start()

	return nil
}

func (myself *consumerBase) CloseAndWait() {
	myself.asyncJobWorker.CloseAndWait()
}
