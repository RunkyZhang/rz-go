package models

import "rz/middleware/notifycenter/enumerations"

type SmsUserCallbackRequestDto struct {
	Message     *SmsMessageDto
	Template    *SmsTemplateDto
	UserMessage *SmsUserMessageDto
}

type MessageStateCallbackRequestDto struct {
	Message      interface{}
	MessageState enumerations.MessageState
}

type MessageExpireCallbackRequestDto struct {
	Message      interface{}
	MessageState enumerations.MessageState
}
