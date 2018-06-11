package models

type QYWeixinMessageDto struct {
	MessageBaseDto

	ToParties []string `json:"toParties"`
	ToTags   []string `json:"toTags"`
}
