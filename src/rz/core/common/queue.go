package common

import (
	"sync"
)

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

func (myself *Queue) Head() (interface{}) {
	myself.lock.Lock()
	defer myself.lock.Unlock()

	if 0 == len(myself.items) {
		return nil
	}

	return myself.items[0]
}

func (myself *Queue) Tail() (interface{}) {
	myself.lock.Lock()
	defer myself.lock.Unlock()

	length := len(myself.items)
	if 0 == length {
		return nil
	}

	return myself.items[length-1]
}
func (myself *Queue) Length() (int) {
	return len(myself.items)
}
