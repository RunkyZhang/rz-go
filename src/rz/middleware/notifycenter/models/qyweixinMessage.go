package models

type QYWeixinMessageDto struct {
	BaseMessageDto

	ToParties []string `json:"toParties"`
	ToTags   []string `json:"toTags"`
}
