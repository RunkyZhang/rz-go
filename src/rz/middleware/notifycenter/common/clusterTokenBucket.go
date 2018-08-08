package common

import (
	"time"
	"sync"
	"fmt"
)

// 48948306/30secs
type ClusterTokenBucket struct {
	intervalSeconds int64
	capacity        int
	redisClient     *RedisClient
	namespace       string
	key             string
	lock            sync.Mutex
	cycleSeconds    int
	slots           map[string]int
	today           time.Time
	availableKey    string
	buffer          int
}

func NewClusterTokenBucket(redisClient *RedisClient, namespace string, key string, intervalSeconds int, capacity int) (*ClusterTokenBucket) {
	Assert.IsNotNilToPanic(redisClient, "redisClient")
	Assert.IsNotBlankToPanic(namespace, "namespace")
	Assert.IsNotBlankToPanic(key, "key")
	Assert.IsTrueToPanic(10 <= intervalSeconds, "10 <= intervalSeconds")
	Assert.IsTrueToPanic(0 < capacity, "0 < capacity")

	clusterTokenBucket := &ClusterTokenBucket{
		redisClient:     redisClient,
		namespace:       namespace,
		key:             key,
		cycleSeconds:    60 * 60 * 24,
		capacity:        capacity,
		intervalSeconds: int64(intervalSeconds),
		slots:           make(map[string]int),
	}
	Assert.IsTrueToPanic(0 == (clusterTokenBucket.cycleSeconds % intervalSeconds), "0 == (clusterTokenBucket.cycleSeconds % intervalSeconds)")
	clusterTokenBucket.buffer = clusterTokenBucket.calculateBuffer()
	clusterTokenBucket.refresh(time.Now())

	return clusterTokenBucket
}

func (myself *ClusterTokenBucket) Capability() (int) {
	return myself.capacity
}

func (myself *ClusterTokenBucket) Available() (int, error) {
	myself.lock.Lock()
	defer myself.lock.Unlock()

	slotKey := myself.calculateSlotKey(time.Now())
	available := myself.slots[slotKey]
	if -1 == available {
		return 0, nil
	} else if 0 < available {
		return available, nil
	}

	return myself.getAvailable(slotKey)
}

func (myself *ClusterTokenBucket) Take(count int, waitingSeconds int) (bool, error) {
	startTimestamp := time.Now().Unix()
	for ; ; {
		ok, err := myself.TryTake(count)
		if nil != err {
			return false, err
		}
		if ok {
			return true, nil
		}

		now := time.Now().Unix()
		leavingSeconds := int64(waitingSeconds) - (now - startTimestamp)
		difference := now - myself.today.Unix()
		nextSeconds := myself.intervalSeconds - (difference % myself.intervalSeconds)

		if leavingSeconds < nextSeconds {
			time.Sleep(time.Duration(leavingSeconds) * time.Second)
			return false, nil
		} else {
			time.Sleep(time.Duration(nextSeconds) * time.Second)
		}
	}
}

func (myself *ClusterTokenBucket) TryTake(count int) (bool, error) {
	err := Assert.IsTrueToError(myself.capacity >= count && 0 < count, "myself.capacity >= count && 0 < count")
	if nil != err {
		return false, err
	}

	myself.lock.Lock()
	defer myself.lock.Unlock()

	now := time.Now()
	if now.Day() != myself.today.Day() {
		myself.refresh(now)
	}

	slotKey := myself.calculateSlotKey(now)
	localAvailable := myself.slots[slotKey]
	if -1 == localAvailable {
		return false, nil
	} else if 0 < localAvailable {
		if localAvailable >= count {
			myself.slots[slotKey] = localAvailable - count
			return true, nil
		}
	}

	// when 0 == localAvailable
	required := myself.buffer
	if count > myself.buffer {
		required = count
	}
	remoteAvailable, err := myself.decrementAvailable(slotKey, required)
	if nil != err {
		return false, err
	}
	if 0 > remoteAvailable {
		leaving := remoteAvailable + required
		if 0 < leaving {
			myself.slots[slotKey] = leaving
		} else {
			myself.slots[slotKey] = -1
			return false, nil
		}
	} else {
		myself.slots[slotKey] = required
	}

	localAvailable = myself.slots[slotKey]
	if localAvailable < count {
		return false, nil
	}
	myself.slots[slotKey] = localAvailable - count

	return true, nil
}

// TODO: rich logic
func (myself *ClusterTokenBucket) calculateBuffer() (int) {
	return myself.capacity / 10
}

func (myself *ClusterTokenBucket) refresh(now time.Time) {
	myself.today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	myself.availableKey = myself.formatAvailableKey()
	myself.initSlots(now)
}

func (myself *ClusterTokenBucket) initSlots(now time.Time) {
	firstSlotStart := myself.calculateSlotStart(now)
	cycleSeconds := int64(myself.cycleSeconds)
	perInitTimes := 100
	index := 0
	secondSlotStart := int64(0)
	for ; firstSlotStart != cycleSeconds; firstSlotStart = firstSlotStart + myself.intervalSeconds {
		key := myself.formatSlotKey(firstSlotStart)
		myself.slots[key] = 0

		if index < perInitTimes {
			_, err := myself.getAvailable(key)
			if nil != err {
				err = myself.supplyAvailable(key)
				if nil != err {
					GetLogging().Error(err, "Failed to init slot(%s)", key)
				}
			}
		} else if index == perInitTimes {
			secondSlotStart = firstSlotStart
		}

		index++
	}

	go func() {
		for ; secondSlotStart != cycleSeconds; secondSlotStart = secondSlotStart + myself.intervalSeconds {
			key := myself.formatSlotKey(secondSlotStart)
			_, err := myself.getAvailable(key)
			if nil != err {
				err = myself.supplyAvailable(key)
				if nil != err {
					GetLogging().Error(err, "Failed to init slot(%s)", key)
				}
			}
		}
	}()
}

func (myself *ClusterTokenBucket) supplyAvailable(slotKey string) (error) {
	return myself.redisClient.HashSet(myself.availableKey, slotKey, Int32ToString(myself.capacity))
}

func (myself *ClusterTokenBucket) getAvailable(slotKey string) (int, error) {
	value, err := myself.redisClient.HashGet(myself.availableKey, slotKey)
	if nil != err {
		return 0, err
	}

	available, err := StringToInt32(value)
	if nil != err {
		return 0, err
	}

	return available, nil
}

func (myself *ClusterTokenBucket) decrementAvailable(slotKey string, count int) (int, error) {
	available, err := myself.redisClient.HashDecrement(myself.availableKey, slotKey, count)

	return int(available), err
}

func (myself *ClusterTokenBucket) calculateSlotStart(now time.Time) (int64) {
	difference := now.Unix() - myself.today.Unix()
	return difference - (difference % myself.intervalSeconds)
}

// slotStart < X <= slotEnd
func (myself *ClusterTokenBucket) calculateSlotKey(now time.Time) (string) {
	slotStart := myself.calculateSlotStart(now)

	return myself.formatSlotKey(slotStart)
}

// slotStart-slotEnd
func (myself *ClusterTokenBucket) formatSlotKey(slotStart int64) (string) {
	return fmt.Sprintf("%d-%d", slotStart, slotStart+myself.intervalSeconds)
}

func (myself *ClusterTokenBucket) formatAvailableKey() (string) {
	return fmt.Sprintf(
		"%s_TokenBucket_Available:%s:%s_%d_%d",
		myself.namespace,
		myself.key,
		myself.today.Format("20060102"),
		myself.intervalSeconds,
		myself.capacity)
}
