package consumers

import (
	"gopkg.in/gomail.v2"
	"net/smtp"
	"errors"
	"encoding/json"
	"time"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
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
	MailMessageConsumer.consumeFunc = MailMessageConsumer.consume
	MailMessageConsumer.handleErrorFunc = MailMessageConsumer.handleError
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

func (mailMessageConsumer *mailMessageConsumer) Send(mailMessageDto *models.MailMessageDto) error {
	message := gomail.NewMessage()
	message.SetHeader("From", mailMessageConsumer.From)
	message.SetHeader("To", mailMessageDto.Tos...)
	message.SetHeader("Subject", mailMessageDto.Subject)
	message.SetBody(mailMessageConsumer.ContentType, mailMessageDto.Content)
	//	m.Attach("/home/Alex/lolcat.jpg")

	return mailMessageConsumer.dialer.DialAndSend(message)
}

func (mailMessageConsumer *mailMessageConsumer) consume(jsonString string) (interface{}, error) {
	mailMessageDto := &models.MailMessageDto{}

	err := json.Unmarshal([]byte(jsonString), mailMessageDto)
	if nil != err {
		return mailMessageDto, nil
	}

	if time.Now().Unix() > mailMessageDto.ExpireTime {
		return mailMessageDto, exceptions.MessageExpire
	}

	//return mailMessageDto, mailMessageConsumer.Send(mailMessageDto)

	return mailMessageDto, nil
}

func (*mailMessageConsumer) handleError(messageDto interface{}, err error) (error) {
	mailMessageDto := messageDto.(*models.MailMessageDto)

	sendChannel, err := enumerations.SendChannelToString(MailMessageConsumer.SendChannel)
	if nil != err {
		return err
	}

	var messageState string
	if err == exceptions.MessageExpire {
		messageState, err = enumerations.MessageStateToString(enumerations.Expire)
	} else {
		messageState, err = enumerations.MessageStateToString(enumerations.Error)
	}
	if nil != err {
		messageState = "Unknown"
	}

	mailMessageDto.States = mailMessageDto.States + ">" + messageState
	mailMessageDto.Exception = err.Error()

	bytes, err := json.Marshal(mailMessageDto)
	if nil != err {
		return err
	}

	return global.GetRedisClient().HashSet(global.RedisKeyMessageValues+sendChannel, mailMessageDto.Id, string(bytes))
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
