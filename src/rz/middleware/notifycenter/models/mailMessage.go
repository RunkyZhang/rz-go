package models

type MailMessageDto struct {
	BaseMessageDto

	Subject string `json:"subject"`
}
