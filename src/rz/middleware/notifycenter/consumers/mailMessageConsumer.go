package consumers

import (
	"gopkg.in/gomail.v2"
	"net/smtp"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/managements"
	"time"
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
	MailMessageConsumer.getMessageFunc = MailMessageConsumer.getMessage
	MailMessageConsumer.sendFunc = MailMessageConsumer.Send
	MailMessageConsumer.poToDtoFunc = MailMessageConsumer.poToDto
	MailMessageConsumer.messageManagementBase = &managements.MailMessageManagement.MessageManagementBase
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

func (myself *mailMessageConsumer) Send(messageDto interface{}) error {
	mailMessageDto, ok := messageDto.(*models.MailMessageDto)
	err := common.Assert.IsTrueToError(ok, "messageDto.(*models.MailMessageDto)")
	if nil != err {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", myself.From)
	message.SetHeader("To", mailMessageDto.Tos...)
	message.SetHeader("Subject", mailMessageDto.Subject)
	message.SetBody(myself.ContentType, mailMessageDto.Content)
	//	m.Attach("/home/Alex/lolcat.jpg")

	return myself.dialer.DialAndSend(message)
}

func (myself *mailMessageConsumer) getMessage(messageId int, date time.Time) (interface{}, *models.PoBase, *models.CallbackBasePo, error) {
	mailMessagePo, err := managements.MailMessageManagement.GetById(messageId, date)
	if nil != err {
		return nil, nil, nil, err
	}

	return mailMessagePo, &mailMessagePo.PoBase, &mailMessagePo.CallbackBasePo, nil
}

func (myself *mailMessageConsumer) poToDto(messagePo interface{}) (interface{}) {
	mailMessagePo, ok := messagePo.(*models.MailMessagePo)
	if !ok {
		return nil
	}

	return models.MailMessagePoToDto(mailMessagePo)
}
