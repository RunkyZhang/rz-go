package main

import (
	"os"
	"os/signal"
	"time"

	"rz/core/common"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/healths"
	"rz/middleware/notifycenter/controllers"
	"rz/middleware/notifycenter/consumers"
)

// http://work.weixin.qq.com/api/doc
// https://cloud.tencent.com/document/product/382/5976
// 202067351   Zgadmin0719   qcloud.com
func main() {
	common.GetLogging().Info(nil, "Starting...")

	start()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill)
	<-exit

	stop()

	common.GetLogging().Info(nil, "Stopped...")
}

func start() {
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
	controllers.SmsProviderController.Enable(controllers.SmsProviderController, true)
	controllers.SystemAliasPermissionController.Enable(controllers.SystemAliasPermissionController, false)
	controllers.SmsUserCallbackController.Enable(controllers.SmsUserCallbackController, false)

	// healths
	global.WebService.RegisterHealthIndicator(&healths.RedisHealthIndicator{})
	global.WebService.RegisterHealthIndicator(&healths.MySQLHealthIndicator{})
	global.WebService.RegisterHealthIndicator(&healths.RuntimeHealthIndicator{})

	// web service
	common.GetLogging().Info(nil, "Web service listening %s ...", global.GetConfig().Web.Listen)
	global.WebService.Start()
}

func stop() {
	// web service
	err := global.WebService.Stop()
	if nil != err {
		common.GetLogging().Error(err, "Failed to shutdown web server")
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
	common.GetLogging().Info(nil, "There are (%d) jobs in [AsyncWorker]. waiting it done...", global.AsyncJobWorker.QueueLength())
	global.AsyncJobWorker.CloseAndWait()
	common.GetLogging().Info(nil, "[AsyncWorker] done...")
}
