package models

type SmsMessageDto struct {
	BaseMessageDto

	NationCode   string   `json:"nationCode"`
	CallbackUrls []string `json:"callbackUrls"`
	CallbackTag  string   `json:"callbackTag"`
}
