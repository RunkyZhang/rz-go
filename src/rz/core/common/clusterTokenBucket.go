package common

import (
	"time"
	"sync"
	"fmt"
)

// 20000~25000/secs
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
	preloadCount    int
	stopped         bool
}

func NewClusterTokenBucket(redisClient *RedisClient, namespace string, key string, intervalSeconds int, capacity int) (*ClusterTokenBucket) {
	Assert.IsTrueToPanic(nil != redisClient, "nil != redisClient")
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
		stopped:         false,
	}
	// min is 1 Minute
	clusterTokenBucket.preloadCount = 6
	Assert.IsTrueToPanic(0 == (clusterTokenBucket.cycleSeconds % intervalSeconds), "0 == (clusterTokenBucket.cycleSeconds % intervalSeconds)")
	clusterTokenBucket.buffer = clusterTokenBucket.calculateBuffer()
	clusterTokenBucket.refresh(time.Now())

	return clusterTokenBucket
}

func (myself *ClusterTokenBucket) Take(count int, waitingSeconds int) (bool, error) {
	startTimestamp := time.Now().Unix()
	todayTimestamp := myself.today.Unix()
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
		difference := now - todayTimestamp
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

	err = Assert.IsTrueToError(false == myself.stopped, "false == myself.stopped")
	if nil != err {
		return false, err
	}

	now := time.Now()
	if now.Day() != myself.today.Day() {
		myself.refresh(now)
	}

	slotKey := myself.calculateSlotKey(myself.today, now)
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
	remoteAvailable, err := myself.decrementAvailable(myself.availableKey, slotKey, required)
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

func (myself *ClusterTokenBucket) PreInitSlots(dateTime time.Time) {
	date := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, time.Local)
	availableKey := myself.formatAvailableKey(date)
	myself.initSlots(availableKey, date, dateTime)
}

func (myself *ClusterTokenBucket) Capability() (int) {
	return myself.capacity
}

func (myself *ClusterTokenBucket) Available() (int, error) {
	// lock for block myself.today and myself.availableKe are difference date
	myself.lock.Lock()
	defer myself.lock.Unlock()

	now := time.Now()
	if now.Day() != myself.today.Day() {
		myself.refresh(now)
	}

	slotKey := myself.calculateSlotKey(myself.today, now)
	localAvailable := myself.slots[slotKey]
	if -1 == localAvailable {
		return 0, nil
	}

	remoteAvailable, err := myself.getAvailable(myself.availableKey, slotKey)
	if nil != err {
		return 0, err
	}

	return localAvailable + remoteAvailable, nil
}

func (myself *ClusterTokenBucket) Stop() (error) {
	myself.lock.Lock()
	defer myself.lock.Unlock()

	myself.stopped = true

	now := time.Now()
	if now.Day() != myself.today.Day() {
		myself.refresh(now)
	}

	slotKey := myself.calculateSlotKey(myself.today, now)
	localAvailable := myself.slots[slotKey]
	if -1 == localAvailable || 0 == localAvailable {
		return nil
	}

	remoteAvailable, err := myself.getAvailable(myself.availableKey, slotKey)
	if nil != err {
		return err
	}

	if 0 > remoteAvailable {
		err = myself.supplyAvailable(myself.availableKey, slotKey, localAvailable)
	} else {
		_, err = myself.incrementAvailable(myself.availableKey, slotKey, localAvailable)
	}

	return err
}

// TODO: rich logic
func (myself *ClusterTokenBucket) calculateBuffer() (int) {
	return myself.capacity / 10
}

func (myself *ClusterTokenBucket) refresh(now time.Time) {
	myself.today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	myself.availableKey = myself.formatAvailableKey(myself.today)
	myself.initSlots(myself.availableKey, myself.today, now)
}

func (myself *ClusterTokenBucket) initSlots(availableKey string, today time.Time, now time.Time) {
	firstSlotStart := myself.calculateSlotStart(today, now)
	cycleSeconds := int64(myself.cycleSeconds)
	index := 0
	secondSlotStart := int64(0)
	for ; firstSlotStart != cycleSeconds; firstSlotStart = firstSlotStart + myself.intervalSeconds {
		key := myself.formatSlotKey(firstSlotStart)
		myself.slots[key] = 0

		if index < myself.preloadCount {
			_, err := myself.getAvailable(availableKey, key)
			if nil != err {
				err = myself.supplyAvailable(availableKey, key, myself.capacity)
				if nil != err {
					GetLogging().Error(err, "Failed to init slot(%s)", key)
				}
			}
		} else if index == myself.preloadCount {
			secondSlotStart = firstSlotStart
		}

		index++
	}

	go func() {
		for ; secondSlotStart != cycleSeconds; secondSlotStart = secondSlotStart + myself.intervalSeconds {
			key := myself.formatSlotKey(secondSlotStart)
			_, err := myself.getAvailable(availableKey, key)
			if nil != err {
				err = myself.supplyAvailable(availableKey, key, myself.capacity)
				if nil != err {
					GetLogging().Error(err, "Failed to init slot(%s)", key)
				}
			}
		}
	}()
}

func (myself *ClusterTokenBucket) supplyAvailable(availableKey string, slotKey string, capacity int) (error) {
	return myself.redisClient.HashSet(availableKey, slotKey, Int32ToString(capacity))
}

func (myself *ClusterTokenBucket) getAvailable(availableKey string, slotKey string) (int, error) {
	value, err := myself.redisClient.HashGet(availableKey, slotKey)
	if nil != err {
		return 0, err
	}

	available, err := StringToInt32(value)
	if nil != err {
		return 0, err
	}

	return available, nil
}

func (myself *ClusterTokenBucket) decrementAvailable(availableKey string, slotKey string, count int) (int, error) {
	available, err := myself.redisClient.HashDecrement(availableKey, slotKey, count)

	return int(available), err
}

func (myself *ClusterTokenBucket) incrementAvailable(availableKey string, slotKey string, count int) (int, error) {
	available, err := myself.redisClient.HashIncrement(availableKey, slotKey, count)

	return int(available), err
}

// slotStart < X <= slotEnd
func (myself *ClusterTokenBucket) calculateSlotKey(today time.Time, now time.Time) (string) {
	slotStart := myself.calculateSlotStart(today, now)

	return myself.formatSlotKey(slotStart)
}

func (myself *ClusterTokenBucket) calculateSlotStart(today time.Time, now time.Time) (int64) {
	difference := now.Unix() - today.Unix()
	return difference - (difference % myself.intervalSeconds)
}

// slotStart-slotEnd
func (myself *ClusterTokenBucket) formatSlotKey(slotStart int64) (string) {
	return fmt.Sprintf("%d-%d", slotStart, slotStart+myself.intervalSeconds)
}

func (myself *ClusterTokenBucket) formatAvailableKey(today time.Time) (string) {
	return fmt.Sprintf(
		"%s_TokenBucket_Available:%s:%s_%d_%d",
		myself.namespace,
		myself.key,
		today.Format("20060102"),
		myself.intervalSeconds,
		myself.capacity)
}
