package models

type SmsMessageDto struct {
	BaseMessageDto

	Parameters         []string `json:"parameters"`
	Sign               string   `json:"sign"`
	TemplateId         int      `json:"templateId"`
	NationCode         string   `json:"nationCode"`
	IdentifyingCode    string   `json:"identifyingCode"`
	ExpireCallbackUrls []string `json:"expireCallbackUrls"`
}
