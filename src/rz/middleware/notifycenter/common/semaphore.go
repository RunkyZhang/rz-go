package common

type Semaphore struct {
	queue      Queue
	isReleased bool
}

func (myself *Semaphore) Wait() {
	if myself.isReleased {
		myself.isReleased = false
		return
	}

	myself.isReleased = false
	channel := make(chan bool, 1)
	myself.queue.Enqueue(channel)
	<-channel
	close(channel)
}

func (myself *Semaphore) Release() {
	myself.isReleased = true

	for ; ; {
		item := myself.queue.Dequeue()
		channel, ok := item.(chan bool)
		if !ok {
			break
		}

		channel <- true
	}
}
