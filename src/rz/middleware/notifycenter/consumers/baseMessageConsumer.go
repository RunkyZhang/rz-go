package consumers

import (
	"time"
	"rz/middleware/notifycenter/enumerations"
	"fmt"
)

type baseMessageConsumer struct {
	SendChannel enumerations.SendChannel
}

func (*baseMessageConsumer) Consume(duration time.Duration) {
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			fmt.Println(time.Now())
			timer.Reset(duration)
		}
	}
}
