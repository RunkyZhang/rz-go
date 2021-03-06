package enumerations

import (
	"rz/middleware/notifycenter/exceptions"
)

type SendChannel int

const (
	Mail SendChannel = iota
	Sms
	QYWeixin
	Weixin
	JPush
	Voice
	SmsCallback
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
	} else if SmsCallback == sendChannel {
		return "SmsCallback", nil
	}

	return "", exceptions.InvalidSendChannel()
}

type ToType int

const (
	Auto ToType = iota
	PhoneNumber
	MailAddress
	UserId
)

type QYWeixinMessageType int

const (
	Text QYWeixinMessageType = iota
	TextCard
)

type MessageState int

const (
	Initial MessageState = iota
	Consuming
	Sent
	FinishedSent
	Error
	ExpireConsuming
	ExpireSent
	ExpireError
)

func MessageStateToString(messageState MessageState) (string) {
	if Initial == messageState {
		return "Initial"
	} else if Consuming == messageState {
		return "Consuming"
	} else if Sent == messageState {
		return "Sent"
	} else if FinishedSent == messageState {
		return "FinishedSent"
	} else if Error == messageState {
		return "Error"
	} else if ExpireConsuming == messageState {
		return "ExpireConsuming"
	} else if ExpireSent == messageState {
		return "ExpireSent"
	} else if ExpireError == messageState {
		return "ExpireError"
	} else {
		return "Unknown"
	}
}

type SmsTemplateType int

const (
	Pattern SmsTemplateType = iota
	IdentifyingCode
)

const (
	SmsContextTypeBusiness      = 1
	SmsContextTypeAdvertisement = 2
)
