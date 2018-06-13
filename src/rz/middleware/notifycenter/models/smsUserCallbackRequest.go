package models

type SmsUserCallbackRequestDto struct {
	Message     *SmsMessageDto
	Template    *SmsTemplateDto
	UserMessage *SmsUserMessageDto
}
