package main

import (
	"fmt"
	"os"
	"os/signal"

	"rz/middleware/notifycenter/controllers"
	"rz/middleware/notifycenter/web"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/repositories"
)

// http://work.weixin.qq.com/api/doc
// https://cloud.tencent.com/document/product/382/5976
// 202067351   Zgadmin0719   qcloud.com
func main() {
	var err error
	//ree := managements.SmsTemplateManagement.Set(11722, 3333, nil, "")
	//fmt.Println(ree)
	//jsonString, ree := global.GetRedisClient().HashGet(global.RedisKeySmsTemplates, common.Int32ToString(117232))
	//smsTemplateDto := &models.SmsTemplateDto{}
	//ree = json.Unmarshal([]byte(""), smsTemplateDto)
	//fmt.Println(jsonString, ree)

	repositories.Init(
		map[string]string{
			"default": "ua_notifycenter:ekIxrgWsJ03u@tcp(10.0.34.44:3306)/notifycenter?parseTime=true",
		})

	//asd, err := repositories.SmsMessageRepository.SelectById(100000002, time.Now())
	//fmt.Println(asd, err)

	//consumers.ConsumerStart()

	fmt.Println("start listening", global.Config.Web.Listen, "...")
	controllers.MessageController.Enable()

	web.Start()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill)
	<-exit

	err = web.Stop()
	if nil != err {
		fmt.Println("Failed to shutdown web server. error: ", err, ".")
	}

	fmt.Println("stoped...")
}
