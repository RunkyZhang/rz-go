package common

import (
	"sync"
	"time"
)

type TokenBucket struct {
	lock      sync.Mutex
	interval  time.Duration
	capacity  int
	available int
	queue     Queue
	ticker    *time.Ticker
}

type channelPack struct {
	channel   chan bool
	count     int
	abandoned bool
}

func NewTokenBucket(interval time.Duration, capacity int) (*TokenBucket) {
	Assert.IsTrueToPanic(0 < interval, "0 < interval")
	Assert.IsTrueToPanic(0 < capacity, "0 < capacity")

	tokenBucket := &TokenBucket{
		interval:  interval,
		capacity:  capacity,
		available: capacity,
		ticker:    time.NewTicker(interval),
	}

	go tokenBucket.supply()

	return tokenBucket
}

func (myself *TokenBucket) Capability() (int) {
	return myself.capacity
}

func (myself *TokenBucket) Available() (int) {
	return myself.available
}

func (myself *TokenBucket) Take(count int, waitingTime time.Duration) (bool) {
	if myself.TryTake(count) {
		return true
	}

	channelPack := &channelPack{
		count:     count,
		abandoned: false,
		channel:   make(chan bool, 1),
	}
	myself.queue.Enqueue(channelPack)

	select {
	case <-channelPack.channel:
		close(channelPack.channel)
		return true
	case <-time.After(waitingTime):
		channelPack.abandoned = true
		return false
	}
}

func (myself *TokenBucket) TryTake(count int) (bool) {
	if myself.capacity < count {
		return false
	}

	myself.lock.Lock()
	defer myself.lock.Unlock()

	if myself.available < count {
		return false
	}

	myself.available -= count
	return true
}

func (myself *TokenBucket) supply() {
	for range myself.ticker.C {
		myself.lock.Lock()

		myself.available = myself.capacity
		for ; ; {
			head := myself.queue.Head()
			channelPack, ok := head.(*channelPack)
			if !ok || nil == channelPack.channel {
				break
			}

			if channelPack.abandoned {
				myself.queue.Dequeue()
				close(channelPack.channel)
			} else {
				if myself.available < channelPack.count {
					// when not enough
					break
				} else {
					myself.queue.Dequeue()
					myself.available -= channelPack.count
					channelPack.channel <- true
				}
			}
		}

		myself.lock.Unlock()
	}
}
