package models

import "rz/middleware/notifycenter/enumerations"

type MessageDto struct {
	Id           string                   `json:"id"`
	Content      string                   `json:"content"`
	SendChannel  enumerations.SendChannel `json:"sendChannel"`
	Tos          []string                 `json:"tos"`
	ToType       enumerations.ToType      `json:"toType"`
	ScheduleTime int64                    `json:"scheduleTime"`
	ExpireTime   int64                    `json:"expireTime"`
	Extra        string                   `json:"extra"`
	SystemAlias  string                   `json:"systemAlias"`
	CreatedTime  int64                    `json:"createdTime"`
}
