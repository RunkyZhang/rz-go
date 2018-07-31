package common

import (
	"time"
	"encoding/json"
	"errors"
	"sync"
)

var (
	clusterTokenBucketKey          = "TokenBucket"
	clusterTokenBucketAvailableKey = clusterTokenBucketKey + "_Available"
)

type clusterTokenBucketPo struct {
	LastSupplyTime  int64  `json:"lastSupplyTime"`
	Capacity        int    `json:"capacity"`
	IntervalSeconds int64  `json:"intervalSeconds"`
	MasterId        string `json:"masterId"`
}

// max 1800token/second
type ClusterTokenBucket struct {
	intervalSeconds int64
	capacity        int
	redisClient     *RedisClient
	namespace       string
	key             string
	lastSupplyTime  int64
	ticker          *time.Ticker
	id              string
	timeoutSeconds  int64
	queue           Queue
	lock            sync.Mutex
}

func NewClusterTokenBucket(redisClient *RedisClient, namespace string, key string, lastSupplyTime int64, intervalSeconds int64, capacity int) (*ClusterTokenBucket) {
	Assert.IsNotNilToPanic(redisClient, "redisClient")
	Assert.IsNotBlankToPanic(namespace, "namespace")
	Assert.IsNotBlankToPanic(key, "key")
	Assert.IsTrueToPanic(10 <= intervalSeconds, "10 <= intervalSeconds")
	Assert.IsTrueToPanic(0 < capacity, "0 < capacity")

	clusterTokenBucket := &ClusterTokenBucket{
		redisClient: redisClient,
		namespace:   namespace,
		key:         key,
		ticker:      time.NewTicker(time.Duration(intervalSeconds) * time.Second),
	}
	ipv4s, err := GetIpV4s()
	if nil != err || 0 == len(ipv4s) {
		panic(errors.New("failed to get ip"))
	}
	clusterTokenBucket.id = "666"
	// depend on time server to adjust every node time
	clusterTokenBucket.timeoutSeconds = 2

	now := time.Now().Unix()
	if now > lastSupplyTime+intervalSeconds {
		lastSupplyTime = now
	}
	refreshed := false
	clusterTokenBucketPo, err := clusterTokenBucket.get()
	if nil != err {
		err = clusterTokenBucket.refresh(lastSupplyTime, intervalSeconds, capacity)
		if nil != err {
			panic(err)
		}
		refreshed = true
	} else {
		overSeconds := now - clusterTokenBucketPo.LastSupplyTime - clusterTokenBucketPo.IntervalSeconds
		if clusterTokenBucket.timeoutSeconds < overSeconds {
			err = clusterTokenBucket.refresh(lastSupplyTime, intervalSeconds, capacity)
			if nil != err {
				panic(err)
			}
			refreshed = true
		}
	}

	if refreshed {
		// block concurrency issue
		time.Sleep(1 * time.Second)
		clusterTokenBucketPo, err = clusterTokenBucket.get()
		if nil != err {
			panic(errors.New("failed to get [TokenBucket] from redis"))
		}
	}

	if clusterTokenBucketPo.IntervalSeconds != intervalSeconds {
		panic(errors.New("[intervalSeconds] conflict when concurrency"))
	}
	if clusterTokenBucketPo.Capacity != capacity {
		panic(errors.New("capacity conflict when concurrency"))
	}
	clusterTokenBucket.lastSupplyTime = clusterTokenBucketPo.LastSupplyTime
	clusterTokenBucket.intervalSeconds = clusterTokenBucketPo.IntervalSeconds
	clusterTokenBucket.capacity = clusterTokenBucketPo.Capacity

	go clusterTokenBucket.supply()

	return clusterTokenBucket
}

func (myself *ClusterTokenBucket) Capability() (int) {
	return myself.capacity
}

func (myself *ClusterTokenBucket) Available() (int, error) {
	return myself.getAvailable()
}

func (myself *ClusterTokenBucket) Take(waitingSeconds int) (bool, error) {
	timestamp := time.Now().Unix()
	for ; ; {
		ok, err := myself.TryTake()
		if nil != err {
			return false, err
		}
		if ok {
			return true, nil
		}

		myself.lock.Lock()
		// need in lock, cos sometime wait out lock
		leavingSeconds := int64(waitingSeconds) - (time.Now().Unix() - timestamp)
		if 0 > leavingSeconds {
			myself.lock.Unlock()
			return false, nil
		}

		channelPack := &channelPack{
			count:     1,
			abandoned: false,
			channel:   make(chan bool, 1),
		}
		myself.queue.Enqueue(channelPack)
		myself.lock.Unlock()

		select {
		case <-channelPack.channel:
			close(channelPack.channel)
		case <-time.After(time.Duration(leavingSeconds) * time.Second):
			channelPack.abandoned = true
			return false, nil
		}
	}
}

func (myself *ClusterTokenBucket) TryTake() (bool, error) {
	myself.lock.Lock()
	defer myself.lock.Unlock()

	available, err := myself.decrementAvailable()
	if nil != err {
		return false, err
	}
	if 0 > available {
		return false, nil
	}

	return true, nil
}

func (myself *ClusterTokenBucket) supply() {
	for range myself.ticker.C {
		myself.lock.Lock()
		clusterTokenBucketPo, err := myself.get()
		if nil != err {
			GetLogging().Error(err, "Failed to get [TokenBucket] from redis")
			myself.lock.Unlock()
			continue
		}
		lastSupplyTime := myself.lastSupplyTime
		myself.lastSupplyTime = clusterTokenBucketPo.LastSupplyTime

		now := time.Now().Unix()
		if myself.id == clusterTokenBucketPo.MasterId { // if self is master, then refresh
			err = myself.refresh(now, myself.intervalSeconds, myself.capacity)
			if nil != err {
				myself.lock.Unlock()
				continue
			}
		} else {
			if lastSupplyTime == myself.lastSupplyTime { // if last supply time is not change, then wait a moment or change master
				overSeconds := now - myself.lastSupplyTime - myself.intervalSeconds
				if myself.timeoutSeconds < overSeconds { // if over time greater than timeout time, then change master
					err = myself.refresh(now, myself.intervalSeconds, myself.capacity)
					if nil != err {
						myself.lock.Unlock()
						continue
					}
				} else {
					waitSeconds := int64(0)
					if 0 > overSeconds {
						waitSeconds = myself.timeoutSeconds - overSeconds
					} else {
						waitSeconds = (-1 * waitSeconds) + myself.timeoutSeconds
					}

					if 0 == waitSeconds {
						waitSeconds = 1
					}
					time.Sleep(time.Duration(waitSeconds) * time.Second)
				}
			}
		}

		for ; ; {
			head := myself.queue.Dequeue()
			channelPack, ok := head.(*channelPack)
			if !ok {
				break
			}
			if nil == channelPack.channel {
				continue
			}

			if channelPack.abandoned {
				close(channelPack.channel)
			} else {
				channelPack.channel <- true
			}
		}
		myself.lock.Unlock()
	}
}

func (myself *ClusterTokenBucket) refresh(lastSupplyTime int64, intervalSeconds int64, capacity int) (error) {
	clusterTokenBucketPo := myself.build(lastSupplyTime, intervalSeconds, capacity)
	err := myself.set(clusterTokenBucketPo)
	if nil != err {
		GetLogging().Error(err, "Failed to get [TokenBucket] from redis")
		return err
	}
	err = myself.supplyAvailable(capacity)
	if nil != err {
		GetLogging().Error(err, "Failed to get [TokenBucket] from redis")
		return err
	}

	return nil
}

func (myself *ClusterTokenBucket) get() (*clusterTokenBucketPo, error) {
	jsonString, err := myself.redisClient.HashGet(myself.namespace+clusterTokenBucketKey, myself.key)
	if nil != err {
		return nil, err
	}

	clusterTokenBucketPo := &clusterTokenBucketPo{}
	err = json.Unmarshal([]byte(jsonString), &clusterTokenBucketPo)
	if nil != err {
		return nil, err
	}

	return clusterTokenBucketPo, nil
}

func (myself *ClusterTokenBucket) set(clusterTokenBucketPo *clusterTokenBucketPo) (error) {
	buffer, err := json.Marshal(clusterTokenBucketPo)
	if nil != err {
		return err
	}

	return myself.redisClient.HashSet(myself.namespace+clusterTokenBucketKey, myself.key, string(buffer))
}

func (myself *ClusterTokenBucket) build(lastSupplyTime int64, intervalSeconds int64, capacity int) (*clusterTokenBucketPo) {
	clusterTokenBucketPo := &clusterTokenBucketPo{
		LastSupplyTime:  lastSupplyTime,
		IntervalSeconds: intervalSeconds,
		Capacity:        capacity,
		MasterId:        myself.id,
	}

	return clusterTokenBucketPo
}

func (myself *ClusterTokenBucket) supplyAvailable(capacity int) (error) {
	return myself.redisClient.HashSet(myself.namespace+clusterTokenBucketAvailableKey, myself.key, Int32ToString(capacity))
}

func (myself *ClusterTokenBucket) getAvailable() (int, error) {
	value, err := myself.redisClient.HashGet(myself.namespace+clusterTokenBucketAvailableKey, myself.key)
	if nil != err {
		return 0, err
	}

	available, err := StringToInt32(value)
	if nil != err {
		return 0, err
	}

	return available, nil
}

func (myself *ClusterTokenBucket) decrementAvailable() (int, error) {
	available, err := myself.redisClient.HashDecrement(myself.namespace+clusterTokenBucketAvailableKey, myself.key, 1)

	return int(available), err
}
