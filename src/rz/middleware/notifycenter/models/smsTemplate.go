package models

type SmsTemplateDto struct {
	Id               int      `json:"id"`
	Extend           int      `json:"extend"`
	UserCallbackUrls []string `json:"userCallbackUrls"`
	Pattern          string   `json:"pattern"`
	Disable          bool     `json:"disable"`
}
