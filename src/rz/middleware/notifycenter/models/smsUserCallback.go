package models

type SmsUserCallbackDto struct {
	Id                   string                               `json:"id"`
	NationCode           string                               `json:"nationCode"`
	PhoneNumber          string                               `json:"phoneNumber"`
	TemplateId           int                                  `json:"templateId"`
	MaxExpireTime        int64                                `json:"maxExpireTime"`
	UserCallbackMessages map[string]SmsUserCallbackMessageDto `json:"userCallbackMessages"`
}
