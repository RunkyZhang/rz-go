package main

import (
	"fmt"
	"os"
	"os/signal"

	"rz/middleware/notifycenter/controllers"
	"rz/middleware/notifycenter/web"
	"rz/middleware/notifycenter/global"
)

// http://work.weixin.qq.com/api/doc
func main() {
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
