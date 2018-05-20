package models

type QYWeixinMessageDto struct {
	MessageDto

	ToParties []string `json:"toParties"`
	ToTags   []string `json:"toTags"`
}
