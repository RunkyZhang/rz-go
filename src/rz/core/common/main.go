package common

import (
	"fmt"
	"time"
	"rz/middleware/notifycenter/models"
	"errors"
)

func main()  {
	//for i := 0; i < 100; i++ {
	//	testSemaphore()
	//}
	//testTokenBucket()
	//testHttpClient()
	//testQueue()
	//testAsyncJobWorker()
	//testAsyncJobTrigger()
	//testClusterTokenBucket()
	//testZooKeeperClient()
}

func testQueue() {
	var queue Queue

	values := [20000000]int{}
	length := len(values)

	go func() {
		for i := 0; i < length/2; i++ {
			queue.Enqueue(i)
		}
		fmt.Println("done Enqueue 1")
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("length", queue.Length())

	go func() {
		for i := length / 2; i < length; i++ {
			queue.Enqueue(i)
		}
		fmt.Println("done Enqueue 2")
	}()

	go func() {
		for ; 0 < queue.Length(); {
			value := queue.Dequeue()
			index, ok := value.(int)
			if !ok {
				continue
			}
			values[index] = 88
		}
		fmt.Println("done Dequeue 1")
	}()

	go func() {
		for ; 0 < queue.Length(); {
			value := queue.Dequeue()
			index, ok := value.(int)
			if !ok {
				continue
			}
			values[index] = 88
		}
		fmt.Println("done Dequeue 2")
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("length", queue.Length())

	go func() {
		for ; 0 < queue.Length(); {
			value := queue.Dequeue()
			index, ok := value.(int)
			if !ok {
				continue
			}
			values[index] = 88
		}
		fmt.Println("done Dequeue 3")
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("length", queue.Length())

	for i := 0; i < len(values); i++ {
		if 88 != values[i] {
			fmt.Println("not 88", i)
			break
		}
	}
}

func testHttpClient() {
	httpClient := NewHttpClient(nil)

	for i := 0; i < 20; i++ {
		go func(value int) {
			start := 100000000 + value*10000
			for j := start; j < start+10000; j++ {
				queryMessageRequestDto := &models.QueryMessagesByIdsRequestDto{}
				queryMessageRequestDto.Ids = []string{"100302541"}
				buffer, err := httpClient.Post("http://notifycenter.cloud.zhaogangrentest.com/message/query-sms", queryMessageRequestDto)

				if nil != err {
					fmt.Println(err)
				}

				if 0 == j%999 {
					fmt.Println(j, string(buffer))
				}
			}
		}(i)
	}

	time.Sleep(10 * time.Minute)
}

func testTokenBucket() {
	tokenBucket := NewTokenBucket(10*time.Second, 3)
	//for i := 0; i < 100; i++ {
	//	//time.Sleep(time.Second)
	//	//fmt.Println(tokenBucket.TryTake(1))
	//	tokenBucket.Take(1, 5*time.Second)
	//}

	tokenBucket.Take(3, 5*time.Second)
	fmt.Println("0------3")
	go func() {
		tokenBucket.Take(2, 1000*time.Second)
		fmt.Println("1------2")
	}()
	time.Sleep(1)
	go func() {
		tokenBucket.Take(2, 1000*time.Second)
		fmt.Println("2------2")
	}()
	time.Sleep(1)
	go func() {
		tokenBucket.Take(1, 1000*time.Second)
		fmt.Println("3------1")
	}()
	time.Sleep(2000 * time.Second)

	timePoint := time.Now().Unix()
	count := 0
	disable := false
	for i := 0; i < 100; i++ {
		go func(index int) {
			if 0 == index%2 {
				for ; ; {
					if tokenBucket.Take(1, 2*time.Second) {
						count += 1
					}
					if disable {
						break
					}
				}
			} else {
				for ; ; {
					if tokenBucket.TryTake(5) {
						count += 5
					}
					if disable {
						break
					}
				}
			}
		}(i)
	}

	fmt.Println("wait 30")
	for ; time.Now().Unix()-timePoint < 30; {
		time.Sleep(1 * time.Millisecond)
	}
	disable = true
	fmt.Println(time.Now().Unix()-timePoint, count)
}

func testSemaphore() {
	semaphore := Semaphore{}

	values := [100]int{}
	for i := 0; i < len(values); i++ {
		if 0 == i%2 {
			go func(index int) {
				semaphore.Release()
				values[index] = index
			}(i)
		} else {
			go func(index int) {
				semaphore.Wait()
				values[index] = -1
			}(i)
		}
	}

	time.Sleep(5 * time.Second)

	semaphore.Release()
	semaphore.Wait()

	for i := 0; i < len(values); i++ {
		if 0 == i%2 {
		} else {
			if -1 != values[i] {
				fmt.Println("not [-1]", i)
				time.Sleep(5 * time.Second)
				if -1 != values[i] {
					fmt.Println("***not [-1] again***", i)
				}
			}
		}
	}
	fmt.Println("done")
}

func testAsyncJobWorker() {
	asyncWorker := NewAsyncJobWorker(5)
	asyncWorker.Start()
	time.Sleep(1 * time.Second)
	for i := 0; i < 10; i++ {
		asyncJob := &AsyncJob{
			Name:      fmt.Sprint(i),
			Type:      "777",
			Parameter: []interface{}{i},
			RunFunc: func(parameter interface{}) error {
				time.Sleep(1 * time.Second)
				return errors.New(fmt.Sprint("test", parameter))
			},
		}
		asyncWorker.Add(asyncJob)
	}

	time.Sleep(10 * time.Second)

	for i := 0; i < 20; i++ {
		asyncJob := &AsyncJob{
			Name:      fmt.Sprint(i),
			Type:      "777",
			Parameter: []interface{}{i},
			RunFunc: func(parameter interface{}) error {
				time.Sleep(5 * time.Second)
				return errors.New(fmt.Sprint("test", parameter))
			},
		}
		asyncWorker.Add(asyncJob)
	}

	fmt.Println(time.Now(), "CloseAndWait")
	asyncWorker.CloseAndWait()
}

func testAsyncJobTrigger() {
	asyncJob := AsyncJob{
		Name: "666",
		Type: "777",
		RunFunc: func(parameter interface{}) error {
			//time.Sleep(5 * time.Second)
			return errors.New(fmt.Sprint("test", parameter))
		},
	}
	asyncJobTrigger := NewAsyncJobTrigger(5, 1*time.Second, asyncJob)
	asyncJobTrigger.Start()

	time.Sleep(100 * time.Second)

	fmt.Println(time.Now(), "CloseAndWait")
	asyncJobTrigger.CloseAndWait()
}

func testClusterTokenBucket() {
	redisClientSettings := RedisClientSettings{
		PoolMaxActive:   10,
		PoolMaxIdle:     1,
		PoolWait:        true,
		PoolIdleTimeout: 180 * time.Second,
		DatabaseId:      0,
		ConnectTimeout:  2000 * time.Second,
		Address:         "10.0.52.105:6379",
		Password:        "",
	}
	redisClient := NewRedisClient(redisClientSettings)

	clusterTokenBucket := NewClusterTokenBucket(redisClient, "Middleware_NotifyCenter", "notifycenter.test", 10, 10)
	//for i := 0; i < 100; i++ {
	//	//time.Sleep(time.Second)
	//	//fmt.Println(clusterTokenBucket.TryTake(1))
	//	clusterTokenBucket.Take(1, 7)
	//}

	startTime := time.Now()
	fmt.Println(startTime)
	count := 0
	disable := false
	for i := 0; i < 100; i++ {
		go func(index int) {
			if 0 == index%2 {
				for ; ; {
					perCount := 1
					ok, _ := clusterTokenBucket.Take(perCount, 2)
					if ok {
						count += perCount
					}
					if disable {
						break
					}
				}
			} else {
				for ; ; {
					perCount := 5
					ok, _ := clusterTokenBucket.TryTake(perCount)
					if ok {
						count += perCount
					}
					if disable {
						break
					}
				}
			}
		}(i)
	}

	seconds := int64(30)
	fmt.Println("wait", seconds)
	for ; time.Now().Unix()-startTime.Unix() < seconds; {
		time.Sleep(1 * time.Millisecond)
	}
	disable = true
	fmt.Println(time.Now().Unix()-startTime.Unix(), count)
}