package channels

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"gopkg.in/gomail.v2"
	"net/smtp"
	"strings"
)

var (
	MailDefaultChannel *mailDefaultChannel
)

func init() {
	MailDefaultChannel = &mailDefaultChannel{
		Host:        global.GetConfig().Mail.Host,
		Port:        global.GetConfig().Mail.Port,
		UserName:    global.GetConfig().Mail.UserName,
		password:    global.GetConfig().Mail.Password,
		From:        global.GetConfig().Mail.From,
		ContentType: global.GetConfig().Mail.ContentType,
	}
	MailDefaultChannel.dialer = gomail.NewDialer(MailDefaultChannel.Host, MailDefaultChannel.Port, MailDefaultChannel.UserName, MailDefaultChannel.password)
	MailDefaultChannel.dialer.Auth = &unencryptedAuth{
		smtp.PlainAuth(
			"",
			MailDefaultChannel.UserName,
			MailDefaultChannel.password,
			MailDefaultChannel.Host,
		),
	}
	MailDefaultChannel.mailDoFunc = MailDefaultChannel.do
	MailDefaultChannel.Id = 0

	MailChannels[MailDefaultChannel.Id] = &MailDefaultChannel.mailChannelBase
}

type mailDefaultChannel struct {
	mailChannelBase

	Host        string
	Port        int
	UserName    string
	password    string
	From        string
	ContentType string

	dialer *gomail.Dialer
}

func (myself *mailDefaultChannel) do(mailMessagePo *models.MailMessagePo) (error) {
	tos := strings.Split(mailMessagePo.Tos, ",")

	message := gomail.NewMessage()
	message.SetHeader("From", myself.From)
	message.SetHeader("To", tos...)
	message.SetHeader("Subject", mailMessagePo.Subject)
	message.SetBody(myself.ContentType, mailMessagePo.Content)
	//	m.Attach("/home/Alex/lolcat.jpg")

	return myself.dialer.DialAndSend(message)
}
