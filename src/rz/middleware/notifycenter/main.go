package main

import (
	"fmt"
	"os"
	"os/signal"

	"rz/middleware/notifycenter/controllers"
	"rz/middleware/notifycenter/web"
)

// http://work.weixin.qq.com/api/doc
func main() {
	fmt.Println("starting...")
	controllers.Controller.Enable()

	web.Start()
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill)
	test := <-exit
	fmt.Println(test)
	web.Stop()

	fmt.Println("stoped...")
}
