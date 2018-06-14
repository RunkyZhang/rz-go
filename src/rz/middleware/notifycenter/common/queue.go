package common

import (
	"sync"
)

func newQueue() (*Queue) {
	queue := &Queue{
		items: make([]interface{}, 0),
	}

	return queue
}

type Queue struct {
	items []interface{}
	lock  sync.Mutex
}

func (myself *Queue) Enqueue(item interface{}) {
	myself.lock.Lock()
	defer myself.lock.Unlock()

	myself.items = append(myself.items, item)
}

func (myself *Queue) Dequeue() (interface{}) {
	myself.lock.Lock()
	defer myself.lock.Unlock()

	if 0 == len(myself.items) {
		return nil
	}

	item := myself.items[0]
	myself.items = myself.items[1:]

	return item
}

func (myself *Queue) Length() (int) {
	return len(myself.items)
}
