package provider

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"gopkg.in/gomail.v2"
	"net/smtp"
	"strings"
)

var (
	MailDefaultProvider *mailDefaultProvider
)

func init() {
	MailDefaultProvider = &mailDefaultProvider{
		Host:        global.GetConfig().Mail.Host,
		Port:        global.GetConfig().Mail.Port,
		UserName:    global.GetConfig().Mail.UserName,
		password:    global.GetConfig().Mail.Password,
		From:        global.GetConfig().Mail.From,
		ContentType: global.GetConfig().Mail.ContentType,
	}
	MailDefaultProvider.dialer = gomail.NewDialer(MailDefaultProvider.Host, MailDefaultProvider.Port, MailDefaultProvider.UserName, MailDefaultProvider.password)
	MailDefaultProvider.dialer.Auth = &unencryptedAuth{
		smtp.PlainAuth(
			"",
			MailDefaultProvider.UserName,
			MailDefaultProvider.password,
			MailDefaultProvider.Host,
		),
	}
	MailDefaultProvider.mailDoFunc = MailDefaultProvider.do
	MailDefaultProvider.Id = "mailDefaultProvider"

	mailProviders[MailDefaultProvider.Id] = &MailDefaultProvider.mailProviderBase
}

type mailDefaultProvider struct {
	mailProviderBase

	Host        string
	Port        int
	UserName    string
	password    string
	From        string
	ContentType string

	dialer *gomail.Dialer
}

func (myself *mailDefaultProvider) do(mailMessagePo *models.MailMessagePo) (error) {
	tos := strings.Split(mailMessagePo.Tos, ",")

	message := gomail.NewMessage()
	message.SetHeader("From", myself.From)
	message.SetHeader("To", tos...)
	message.SetHeader("Subject", mailMessagePo.Subject)
	message.SetBody(myself.ContentType, mailMessagePo.Content)
	//	m.Attach("/home/Alex/lolcat.jpg")

	return myself.dialer.DialAndSend(message)
}
