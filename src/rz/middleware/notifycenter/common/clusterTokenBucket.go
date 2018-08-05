package common

import (
	"time"
	"encoding/json"
	"sync"
	"fmt"
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
	id              string
	queue           Queue
	lock            sync.Mutex
	cycleSeconds    int
	slots           map[string]int
	today           time.Time
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
	}
	Assert.IsTrueToPanic(0 == (clusterTokenBucket.cycleSeconds % intervalSeconds), "0 == (clusterTokenBucket.cycleSeconds % intervalSeconds)")
	clusterTokenBucket.initSlots()

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

func (myself *ClusterTokenBucket) initSlots() {
	for i := int64(0); i != int64(myself.cycleSeconds); i = i + myself.intervalSeconds {
		myself.slots[myself.formatSlotKey(int64(i))] = myself.capacity
	}
}

// slotStart < X <= slotEnd
func (myself *ClusterTokenBucket) calculateSlotKey() (string) {
	now := time.Now()
	if now.Day() != myself.today.Day() {
		myself.lock.Lock()
		defer myself.lock.Unlock()
		if now.Day() != myself.today.Day() {
			myself.initSlots()
		}
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	difference := now.Unix() - today.Unix()
	slotStart := difference - (difference % myself.intervalSeconds)

	return myself.formatSlotKey(slotStart)
}

// slotStart-slotEnd
func (myself *ClusterTokenBucket) formatSlotKey(slotStart int64) (string) {
	return fmt.Sprintf("%d-%d", slotStart, slotStart+myself.intervalSeconds)
}
