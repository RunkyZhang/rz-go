package main

import (
	"fmt"
	"os"
	"os/signal"

	"rz/middleware/notifycenter/controllers"
	"rz/middleware/notifycenter/web"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/consumers"
	"time"
)

// http://work.weixin.qq.com/api/doc
// https://cloud.tencent.com/document/product/382/5976
func main() {
	fmt.Println(global.GetRedisClient().Ping())

	go consumers.SmsMessageConsumer.Consume(5 * time.Second)

	fmt.Println("start listening", global.Config.Web.Listen, "...")
	controllers.MessageController.Enable()

	web.Start()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill)
	<-exit

	err := web.Stop()
	if nil != err {
		fmt.Println("Failed to shutdown web server. error:", err, ".")
	}

	fmt.Println("stoped...")
}
