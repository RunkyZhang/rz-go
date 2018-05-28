package models

import "rz/middleware/notifycenter/enumerations"

type MessageDto struct {
	Content      string              `json:"content"`
	Tos          []string            `json:"tos"`
	ToType       enumerations.ToType `json:"toType"`
	ScheduleTime int64               `json:"scheduleTime"`
	ExpireTime   int64               `json:"expireTime"`
	Extra        string              `json:"extra"`
	SystemAlias  string              `json:"systemAlias"`

	Id          string                   `json:"id"`
	SendChannel enumerations.SendChannel `json:"sendChannel"`
	CreatedTime int64                    `json:"createdTime"`
	States      string                   `json:"state"`
	Exception   string                   `json:"exception"`
}
