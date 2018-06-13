package main

import (
	"fmt"
	"os"
	"os/signal"

	"rz/middleware/notifycenter/controllers"
	"rz/middleware/notifycenter/web"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/consumers"
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

	//err = exceptions.DtoNull().AttachMessage("asdasdasd")
	//fmt.Printf("failed to get message ids. error: %s", err)

	repositories.Init(
		map[string]string{
			"default": "ua_notifycenter:ekIxrgWsJ03u@tcp(10.0.34.44:3306)/notifycenter?parseTime=true",
		})

	//smsTemplatePo := &models.SmsTemplatePo{
	//	Id:               11722,
	//	Extend:           3333,
	//	UserCallbackUrls: "111@111,222@222",
	//}
	//asd, err := managements.SmsTemplateManagement.GetByTemplateId(11722)
	//fmt.Println(asd, err)

	//var mailMessagePos []models.MailMessagePo
	//fmt.Println(mailMessagePos)
	//smsMessagePo, err := repositories.SmsMessageRepository.SelectById(100000002, time.Now())
	//repositories.SmsMessageRepository.UpdateById(smsMessagePo.Id, "Test", false, "", time.Now())
	//fmt.Println(smsMessagePo, err)

	consumers.ConsumerStart()

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
