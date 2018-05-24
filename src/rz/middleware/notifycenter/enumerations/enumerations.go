package enumerations

import "rz/middleware/notifycenter/exceptions"

type SendChannel int

const (
	Email    SendChannel = iota
	Sms
	QYWeixin
	Weixin
	JPush
	Voice
)

func SendChannelToString(sendChannel SendChannel) (string, error) {
	if Email == sendChannel {
		return "Email", nil
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
	Phone
	UserId
	Mail
)

type QYWeixinMessageType int

const (
	Text     QYWeixinMessageType = iota
	TextCard
)
