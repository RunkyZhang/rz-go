package models

type SmsMessageDto struct {
	BaseMessageDto

	NationCode   int      `json:"nation_code"`
	CallbackUrls []string `json:"callbackUrls"`
	CallbackTag  string   `json:"callbackTag"`
}
