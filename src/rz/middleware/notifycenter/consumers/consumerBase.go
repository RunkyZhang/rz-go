package consumers

import (
	"time"
	"runtime"
	"rz/middleware/notifycenter/common"
)

type consumerBase struct {
	name            string
	asyncJobTrigger *common.AsyncJobTrigger
	runFunc         common.RunFunc
}

func (myself *consumerBase) Start(interval time.Duration) (error) {
	asyncJob := common.AsyncJob{
		Name:    myself.name,
		Type:    "Consumer",
		RunFunc: myself.runFunc,
	}
	myself.asyncJobTrigger = common.NewAsyncJobTrigger(runtime.NumCPU(), interval, asyncJob)
	myself.asyncJobTrigger.Start()

	return nil
}

func (myself *consumerBase) CloseAndWait() {
	myself.asyncJobTrigger.CloseAndWait()
}
