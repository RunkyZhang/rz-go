package consumers

import (
	"gopkg.in/gomail.v2"
	"net/smtp"
	"encoding/json"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/common"
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

	var err error
	MailMessageConsumer.SendChannel = enumerations.Mail
	MailMessageConsumer.keySuffix, err = enumerations.SendChannelToString(MailMessageConsumer.SendChannel)
	common.Assert.IsNilError(err, "")
	MailMessageConsumer.convertFunc = MailMessageConsumer.convert
	MailMessageConsumer.sendFunc = MailMessageConsumer.Send
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
	messageConsumerBase

	Host        string
	Port        int
	UserName    string
	password    string
	From        string
	ContentType string

	dialer *gomail.Dialer
}

func (mailMessageConsumer *mailMessageConsumer) Send(messageDto interface{}) error {
	mailMessageDto := messageDto.(*models.MailMessageDto)

	return nil

	message := gomail.NewMessage()
	message.SetHeader("From", mailMessageConsumer.From)
	message.SetHeader("To", mailMessageDto.Tos...)
	message.SetHeader("Subject", mailMessageDto.Subject)
	message.SetBody(mailMessageConsumer.ContentType, mailMessageDto.Content)
	//	m.Attach("/home/Alex/lolcat.jpg")

	return mailMessageConsumer.dialer.DialAndSend(message)
}

func (mailMessageConsumer *mailMessageConsumer) convert(jsonString string) (interface{}, *models.MessageBaseDto, error) {
	mailMessageDto := &models.MailMessageDto{}

	err := json.Unmarshal([]byte(jsonString), mailMessageDto)
	if nil != err {
		return nil, nil, err
	}

	return mailMessageDto, &mailMessageDto.MessageBaseDto, nil
}
