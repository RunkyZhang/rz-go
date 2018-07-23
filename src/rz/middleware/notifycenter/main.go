package main

import (
	"os"
	"os/signal"
	"time"
	"fmt"
	"strings"

	"rz/middleware/notifycenter/controllers"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/consumers"
	"rz/middleware/notifycenter/healths"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/models"
)

// http://work.weixin.qq.com/api/doc
// https://cloud.tencent.com/document/product/382/5976
// 202067351   Zgadmin0719   qcloud.com
func main() {
	common.GetLogging().Info(nil, "starting...")

	start()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill)
	<-exit

	stop()

	common.GetLogging().Info(nil, "stopped...")
}

func start() {
	// repositories
	common.SetConnectionStrings(global.GetConfig().Databases)

	// asyncWorker
	global.AsyncWorker.Start()

	// consumers
	duration := time.Duration(global.GetConfig().ConsumingInterval) * time.Second
	consumers.SmsMessageConsumer.Start(duration)
	consumers.MailMessageConsumer.Start(duration)
	consumers.SmsUserMessageConsumer.Start(duration)

	// controllers
	controllers.MessageController.Enable(controllers.MessageController, true)
	controllers.SmsTemplateController.Enable(controllers.SmsTemplateController, true)
	controllers.SmsUserCallbackController.Enable(controllers.SmsUserCallbackController, false)

	// healths
	global.WebService.RegisterHealthIndicator(&healths.RedisHealthIndicator{})
	global.WebService.RegisterHealthIndicator(&healths.MySQLHealthIndicator{})
	global.WebService.RegisterHealthIndicator(&healths.RuntimeHealthIndicator{})

	// web service
	common.GetLogging().Info(nil, "web service listening %s ...", global.GetConfig().Web.Listen)
	global.WebService.Start()
}

func stop() {
	// web service
	err := global.WebService.Stop()
	if nil != err {
		common.GetLogging().Error(err, "failed to shutdown web server")
	}

	// consumers
	common.GetLogging().Info(nil, "[SmsMessageConsumer] closing...")
	consumers.SmsMessageConsumer.CloseAndWait()
	common.GetLogging().Info(nil, "[SmsMessageConsumer] closed...")
	common.GetLogging().Info(nil, "[MailMessageConsumer] closing...")
	consumers.MailMessageConsumer.CloseAndWait()
	common.GetLogging().Info(nil, "[MailMessageConsumer] closed...")
	common.GetLogging().Info(nil, "[SmsUserMessageConsumer] closing...")
	consumers.SmsUserMessageConsumer.CloseAndWait()
	common.GetLogging().Info(nil, "[SmsUserMessageConsumer] closed...")

	// AsyncWorker
	common.GetLogging().Info(nil, "there are (%d) jobs in [AsyncWorker]. waiting it done...", global.AsyncWorker.QueueLength())
	global.AsyncWorker.CloseAndWait()
	common.GetLogging().Info(nil, "[AsyncWorker] done...")
}

//func test() {
//var asd interface{}
//asd = nil
//aj, ok := asd.(*common.AsyncJob)
//fmt.Println(aj, ok)
//ree := managements.SmsTemplateManagement.Set(11722, 3333, nil, "")
//fmt.Println(ree)
//jsonString, ree := global.GetRedisClient().HashGet(global.RedisKeySmsTemplates, common.Int32ToString(117232))
//smsTemplateDto := &models.SmsTemplateDto{}
//ree = json.Unmarshal([]byte(""), smsTemplateDto)
//fmt.Println(jsonString, ree)

//err = exceptions.DtoNull().AttachMessage("asdasdasd")
//global.GetLogging().Info(nil, "failed to get message ids. error: %s", err)

//asyncJob := &common.AsyncJob{
//	Name: "666",
//	Type: "777",
//	RunFunc: func(parameter interface{}) error {
//		time.Sleep(5 * time.Second)
//		//fmt.Println(time.Now())
//
//		return errors.New("test")
//	},
//}

//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//global.AsyncWorker.Add(asyncJob)
//}

func testQueue() {
	var queue common.Queue

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
	httpClient := common.NewHttpClient(nil)

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
	tokenBucket := common.NewTokenBucket(5*time.Second, 20000)
	timePoint := time.Now().Unix()
	//for i := 0; i < 100; i++ {
	//	//time.Sleep(time.Second)
	//	//fmt.Println(tokenBucket.TryTake(1))
	//	fmt.Println(tokenBucket.Take(1, 5*time.Second))
	//}

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
	semaphore := common.Semaphore{}

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

func test() {
	asdlist := []int{1, 2, 8, 4, 5}
	fmt.Println(asdlist[len(asdlist)-1])

	channel := make(chan bool, 1)
	channel <- true
	fmt.Println(<-channel)

	var t1, t2 time.Time
	fmt.Println(t1 == t2)

	var asds []string
	fmt.Println("" == strings.Join(asds, ","))

	asd := strings.Split("%s%s您的提货信息为：提货方式：；货物信息：%s ；共%s；提货函%s已经贵司确认，确认请回复。", "%s")
	fmt.Println(len(asd))

	parameters := []string{}

	var args []interface{}
	for _, parameter := range parameters {
		args = append(args, parameter)
	}

	value := fmt.Sprintf("111%s222%s333", args...)
	fmt.Println(value)

	var values []int64
	values = append(values, 123)
	values = append(values, 121)
	values = append(values, 123)
	values = append(values, 103)

	common.SortReverseIntSlice(values)
	fmt.Println(values)
}
