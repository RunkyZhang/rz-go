package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"rz/middleware/notifycenter/controllers"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/consumers"
	"rz/middleware/notifycenter/healths"
	"rz/middleware/notifycenter/common"
)

// http://work.weixin.qq.com/api/doc
// https://cloud.tencent.com/document/product/382/5976
// 202067351   Zgadmin0719   qcloud.com
func main() {
	fmt.Printf("starting...\n")

	start()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill)
	<-exit

	stop()

	fmt.Printf("stoped...\n")
}

func start() {
	// repositories
	common.InitDatabases(
		map[string]string{
			"default": "ua_notifycenter:ekIxrgWsJ03u@tcp(10.0.34.44:3306)/notifycenter?parseTime=true",
		})

	// asyncWorker
	global.AsyncWorker.Start()

	// consumers
	consumers.SmsMessageConsumer.Start(5 * time.Second)
	consumers.MailMessageConsumer.Start(5 * time.Second)
	consumers.SmsUserMessageConsumer.Start(5 * time.Second)

	// controllers
	controllers.MessageController.Enable(controllers.MessageController, true)
	controllers.SmsUserCallbackController.Enable(controllers.SmsUserCallbackController, false)

	// healths
	redisHealthIndicator, err := healths.NewRedisHealthIndicator(global.RedisClient)
	if nil == err {
		global.WebService.RegisterHealthIndicator(redisHealthIndicator)
	}
	mysqlHealthIndicator, err := healths.NewMySQLHealthIndicator(common.Databases)
	if nil == err {
		global.WebService.RegisterHealthIndicator(mysqlHealthIndicator)
	}

	// web service
	fmt.Printf("web service listening %s ...\n", global.Config.Web.Listen)
	global.WebService.Start()
}

func stop() {
	// web service
	err := global.WebService.Stop()
	if nil != err {
		fmt.Printf("failed to shutdown web server. error: %s\n", err.Error())
	}

	// consumers
	fmt.Printf("[SmsMessageConsumer] closing...")
	consumers.SmsMessageConsumer.CloseAndWait()
	fmt.Printf("[SmsMessageConsumer] closed...")
	fmt.Printf("[MailMessageConsumer] closing...")
	consumers.MailMessageConsumer.CloseAndWait()
	fmt.Printf("[MailMessageConsumer] closed...")
	fmt.Printf("[SmsUserMessageConsumer] closing...")
	consumers.SmsUserMessageConsumer.CloseAndWait()
	fmt.Printf("[SmsUserMessageConsumer] closed...")

	// AsyncWorker
	fmt.Printf("there are (%d) jobs in [AsyncWorker]. waiting it done...\n", global.AsyncWorker.QueueLength())
	global.AsyncWorker.CloseAndWait()
	fmt.Printf("[AsyncWorker] done...\n")
}

func test() {
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
	//fmt.Printf("failed to get message ids. error: %s", err)

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
}
