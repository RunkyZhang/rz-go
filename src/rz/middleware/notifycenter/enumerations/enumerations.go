package enumerations

import "rz/middleware/notifycenter/exceptions"

type SendChannel int

const (
	Mail     SendChannel = iota
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
	Auto        ToType = iota
	PhoneNumber
	UserId
	MailAddress
)

type QYWeixinMessageType int

const (
	Text     QYWeixinMessageType = iota
	TextCard
)

type MessageState int

const (
	Initial      MessageState = iota
	Consuming
	Sent
	AppCallback
	UserCallback
	BothCallback
	Error
)

func MessageStateToString(messageState MessageState) (string, error) {
	if Initial == messageState {
		return "Initial", nil
	} else if Consuming == messageState {
		return "Consuming", nil
	} else if Sent == messageState {
		return "Sent", nil
	} else if AppCallback == messageState {
		return "AppCallback", nil
	} else if UserCallback == messageState {
		return "UserCallback", nil
	} else if BothCallback == messageState {
		return "BothCallback", nil
	} else if Error == messageState {
		return "Error", nil
	}

	return "", exceptions.InvalidMessageState
}
