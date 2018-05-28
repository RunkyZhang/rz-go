package consumers

import (
	"gopkg.in/gomail.v2"
	"net/smtp"
	"rz/middleware/notifycenter/global"
	"git.zhaogangren.com/cloud/cloud.appgov.notifycenter.service/model"
	"rz/middleware/notifycenter/enumerations"
	"errors"
)

var (
	MailMessageConsumer *mailMessageConsumer
)

func init() {
	MailMessageConsumer = &mailMessageConsumer{
		Host:        global.Config.Mail.Host,
		Port:        global.Config.Mail.Port,
		UserName:    global.Config.Mail.UserName,
		password:    global.Config.Mail.Password,
		From:        global.Config.Mail.From,
		ContentType: global.Config.Mail.ContentType,
	}

	MailMessageConsumer.SendChannel = enumerations.Mail
	MailMessageConsumer.dialer = gomail.NewDialer(MailMessageConsumer.Host, MailMessageConsumer.Port, MailMessageConsumer.UserName, MailMessageConsumer.password)
	MailMessageConsumer.dialer.Auth = &unencryptedAuth{
		smtp.PlainAuth(
			"",
			MailMessageConsumer.UserName,
			MailMessageConsumer.password,
			MailMessageConsumer.Host,
		),
	}
}

type mailMessageConsumer struct {
	baseMessageConsumer

	Host        string
	Port        int
	UserName    string
	password    string
	From        string
	ContentType string

	dialer *gomail.Dialer
}

func (mailMessageConsumer *mailMessageConsumer) Send(mailMessageDto *model.MailMessageDto) error {
	message := gomail.NewMessage()
	message.SetHeader("From", mailMessageConsumer.From)
	message.SetHeader("To", mailMessageDto.Tos...)
	message.SetHeader("Subject", mailMessageDto.Subject)
	message.SetBody(mailMessageConsumer.ContentType, mailMessageDto.Content)
	//	m.Attach("/home/Alex/lolcat.jpg")

	return mailMessageConsumer.dialer.DialAndSend(message)
}

type unencryptedAuth struct {
	smtp.Auth
}

func (unencryptedAuth *unencryptedAuth) Start(serverInfo *smtp.ServerInfo) (string, []byte, error) {
	(*serverInfo).TLS = true
	return unencryptedAuth.Auth.Start(serverInfo)
}

func (*unencryptedAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// We've already sent everything.
		return nil, errors.New("unexpected server challenge")
	}

	return nil, nil
}
