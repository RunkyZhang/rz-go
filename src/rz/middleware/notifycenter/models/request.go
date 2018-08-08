package models

import "rz/middleware/notifycenter/enumerations"

type SmsUserCallbackRequestDto struct {
	Message     *SmsMessageDto     `json:"message"`
	Template    *SmsTemplateDto    `json:"template"`
	UserMessage *SmsUserMessageDto `json:"userMessage"`
}

type MessageStateCallbackRequestDto struct {
	Message      interface{}               `json:"message"`
	MessageState enumerations.MessageState `json:"messageState"`
}

type MessageExpireCallbackRequestDto struct {
	Message interface{} `json:"message"`
}

type QueryMessagesRequestDto struct {
	Id          int64  `json:"id,string"`
	SystemAlias string `json:"systemAlias"`
	Year        int    `json:"year"`
}

type QueryMessagesByIdsRequestDto struct {
	Ids []string `json:"ids"`
}

type DisableMessageRequestDto struct {
	Id          int64  `json:"id,string"`
	SystemAlias string `json:"systemAlias"`
}

type QuerySmsUserMessagesRequestDto struct {
	SmsMessageId int64  `json:"smsMessageId,string"`
	Content      string `json:"content"`
	NationCode   string `json:"nationCode"`
	PhoneNumber  string `json:"phoneNumber"`
	TemplateId   int    `json:"templateId"`
	Year         int    `json:"year"`
}

type ModifySystemAliasPermissionRequestDto struct {
	SystemAlias        string `json:"systemAlias"`
	SmsPermission      *int   `json:"smsPermission,omitempty"`
	MailPermission     *int   `json:"mailPermission,omitempty"`
	SmsDayFrequency    *int   `json:"smsDayFrequency,omitempty"`
	SmsHourFrequency   *int   `json:"smsHourFrequency,omitempty"`
	SmsMinuteFrequency *int   `json:"smsMinuteFrequency,omitempty"`
}

type TakeTokenRequestDto struct {
	SystemAlias     string `json:"systemAlias"`
	IntervalSeconds int    `json:"intervalSeconds"`
	Capacity        int    `json:"capacity"`
}
