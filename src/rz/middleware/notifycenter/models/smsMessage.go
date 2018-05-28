package models

type SmsMessageDto struct {
	MessageDto

	NationCode   int      `json:"nation_code"`
	CallbackUrls []string `json:"callbackUrls"`
}
