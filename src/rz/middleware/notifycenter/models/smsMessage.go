package models

type SmsMessageDto struct {
	MessageDto

	CallbackUrls []string `json:"callbackUrls"`
}
