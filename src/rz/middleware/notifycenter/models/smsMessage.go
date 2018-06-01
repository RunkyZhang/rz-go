package models

type SmsMessageDto struct {
	BaseMessageDto

	Sign         string   `json:"sign"`
	TemplateId   int      `json:"templateId"`
	NationCode   string   `json:"nationCode"`
	CallbackUrls []string `json:"callbackUrls"`
	Extend       int      `json:"extend"`
}
