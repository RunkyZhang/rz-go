package models

import "rz/middleware/notifycenter/enumerations"

type BaseMessageDto struct {
	Content              string              `json:"content"`
	Tos                  []string            `json:"tos"`
	ToType               enumerations.ToType `json:"toType"`
	ScheduleTime         int64               `json:"scheduleTime"`
	ExpireTime           int64               `json:"expireTime"`
	Extra                string              `json:"extra"`
	SystemAlias          string              `json:"systemAlias"`
	FinishedCallbackUrls []string            `json:"finishedCallbackUrls"`

	Id           string                   `json:"id"`
	SendChannel  enumerations.SendChannel `json:"sendChannel"`
	CreatedTime  int64                    `json:"createdTime"`
	Finished     bool                     `json:"finished"`
	FinishedTime int64                    `json:"finishedTime"`
	States       string                   `json:"state"`
	ErrorMessage string                   `json:"errorMessage"`
}
