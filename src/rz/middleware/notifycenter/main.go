package main

import (
	"time"
	"fmt"
	"strings"
		"rz/middleware/notifycenter/controllers"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/consumers"
	"rz/middleware/notifycenter/healths"
	"rz/core/common"
		"os"
	"os/signal"
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
	global.AsyncJobWorker.Start()

	// consumers
	duration := time.Duration(global.GetConfig().ConsumingInterval) * time.Second
	consumers.SmsMessageConsumer.Start(duration)
	consumers.MailMessageConsumer.Start(duration)
	consumers.SmsUserMessageConsumer.Start(duration)

	// controllers
	controllers.MessageController.Enable(controllers.MessageController, true)
	controllers.SmsTemplateController.Enable(controllers.SmsTemplateController, true)
	controllers.SmsUserCallbackController.Enable(controllers.SmsUserCallbackController, false)
	controllers.SystemAliasPermissionController.Enable(controllers.SystemAliasPermissionController, false)

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
	common.GetLogging().Info(nil, "there are (%d) jobs in [AsyncWorker]. waiting it done...", global.AsyncJobWorker.QueueLength())
	global.AsyncJobWorker.CloseAndWait()
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
