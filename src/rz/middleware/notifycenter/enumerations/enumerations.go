package enumerations

import "rz/middleware/notifycenter/exceptions"

type SendChannel int

const (
	Mail    SendChannel = iota
	Sms
	QYWeixin
	Weixin
	JPush
	Voice
)

func SendChannelToString(sendChannel SendChannel) (string, error) {
	if Mail == sendChannel {
		return "Mail", nil
	} else if Sms == sendChannel {
		return "Sms", nil
	} else if QYWeixin == sendChannel {
		return "QYWeixin", nil
	} else if Weixin == sendChannel {
		return "Weixin", nil
	} else if JPush == sendChannel {
		return "JPush", nil
	} else if Voice == sendChannel {
		return "Voice", nil
	}

	return "", exceptions.InvalidSendChannel
}

type ToType int

const (
	Auto   ToType = iota
	PhoneNumber
	UserId
	MailAddress
)

type QYWeixinMessageType int

const (
	Text     QYWeixinMessageType = iota
	TextCard
)
