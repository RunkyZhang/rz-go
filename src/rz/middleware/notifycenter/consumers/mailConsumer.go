package consumers

import (
	"gopkg.in/gomail.v2"
	"net/smtp"
	"rz/middleware/notifycenter/global"
	"git.zhaogangren.com/cloud/cloud.appgov.notifycenter.service/model"
)

var (
	MailConsumer *mailConsumer
)

func init() {
	MailConsumer = &mailConsumer{
		Host:        global.Config.Mail.Host,
		Port:        global.Config.Mail.Port,
		UserName:    global.Config.Mail.UserName,
		password:    global.Config.Mail.Password,
		From:        global.Config.Mail.From,
		ContentType: global.Config.Mail.ContentType,
	}

	MailConsumer.dialer = gomail.NewDialer(MailConsumer.Host, MailConsumer.Port, MailConsumer.UserName, MailConsumer.password)
	MailConsumer.dialer.Auth = unencryptedAuth{
		smtp.PlainAuth(
			"",
			MailConsumer.UserName,
			MailConsumer.password,
			MailConsumer.Host,
		),
	}
}

type mailConsumer struct {
	Host        string
	Port        int
	UserName    string
	password    string
	From        string
	ContentType string

	dialer *gomail.Dialer
}

func (mailConsumer *mailConsumer) Send(mailMessageDto *model.MailMessageDto) error {
	message := gomail.NewMessage()
	message.SetHeader("From", mailConsumer.From)
	message.SetHeader("To", mailMessageDto.Tos...)
	message.SetHeader("Subject", mailMessageDto.Subject)
	message.SetBody(mailConsumer.ContentType, mailMessageDto.Content)
	//	m.Attach("/home/Alex/lolcat.jpg")

	return mailConsumer.dialer.DialAndSend(message)
}

type unencryptedAuth struct {
	smtp.Auth
}

func (unencryptedAuth *unencryptedAuth) Start(serverInfo *smtp.ServerInfo) (string, []byte, error) {
	(*serverInfo).TLS = true
	return unencryptedAuth.Auth.Start(serverInfo)
}
