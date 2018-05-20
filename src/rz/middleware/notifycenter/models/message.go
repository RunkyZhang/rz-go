package models

type MessageDto struct {
	//Id string `json:"id"`
	Content string `json:"content"`
	SendChannel SendChannel `json:"sendChannel"`
	Tos []string `json:"tos"`
	ToType ToType `json:"toType"`
	ScheduleTime int64 `json:"scheduleTime"`
	ExpireTime int64 `json:"expireTime"`
	Extra string `json:"extra"`
	SystemAlias string `json:"systemAlias"`
}