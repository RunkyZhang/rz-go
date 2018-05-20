package models

type MailMessageDto struct {
	MessageDto

	Subject string `json:"subject"`
}
