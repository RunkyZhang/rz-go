package models

import (
	"rz/middleware/notifycenter/enumerations"
	"time"
	"strings"
)

type MessageBaseDto struct {
	Content              string              `json:"content"`
	Tos                  []string            `json:"tos"`
	ToType               enumerations.ToType `json:"toType"`
	ScheduleTime         int64               `json:"scheduleTime"`
	ExpireTime           int64               `json:"expireTime"`
	Extra                string              `json:"extra"`
	SystemAlias          string              `json:"systemAlias"`
	FinishedCallbackUrls []string            `json:"finishedCallbackUrls"`

	Id            int                      `json:"id"`
	SendChannel   enumerations.SendChannel `json:"sendChannel"`
	Finished      bool                     `json:"finished"`
	FinishedTime  time.Time                `json:"finishedTime"`
	States        string                   `json:"states"`
	ErrorMessages string                   `json:"errorMessages"`
	CreatedTime   int64                    `json:"createdTime"`
	UpdatedTime   int64                    `json:"updatedTime"`
}

type MessageBasePo struct {
	PoBase

	Content              string              `gorm:"column:content"`
	Tos                  string              `gorm:"column:tos"`
	ToType               enumerations.ToType `gorm:"column:toType"`
	ScheduleTime         time.Time           `gorm:"column:scheduleTime"`
	ExpireTime           time.Time           `gorm:"column:expireTime"`
	Extra                string              `gorm:"column:extra"`
	SystemAlias          string              `gorm:"column:systemAlias"`
	FinishedCallbackUrls string              `gorm:"column:finishedCallbackUrls"`

	Id            int                      `gorm:"column:id;primary_key;auto_increment"`
	SendChannel   enumerations.SendChannel `gorm:"column:sendChannel"`
	Finished      bool                     `gorm:"column:finished"`
	FinishedTime  time.Time                `gorm:"column:finishedTime"`
	States        string                   `gorm:"column:states"`
	ErrorMessages string                   `gorm:"column:errorMessages"`
}

func MessageBaseDtoToPo(messageBaseDto *MessageBaseDto) (*MessageBasePo) {
	messageBasePo := &MessageBasePo{}
	messageBasePo.Content = messageBaseDto.Content
	messageBasePo.Tos = strings.Join(messageBaseDto.Tos, ",")
	messageBasePo.ToType = messageBaseDto.ToType
	messageBasePo.ScheduleTime = time.Unix(messageBaseDto.ScheduleTime, 0)
	messageBasePo.ExpireTime = time.Unix(messageBaseDto.ExpireTime, 0)
	messageBasePo.Extra = messageBaseDto.Extra
	messageBasePo.SystemAlias = messageBaseDto.SystemAlias
	messageBasePo.FinishedCallbackUrls = strings.Join(messageBaseDto.FinishedCallbackUrls, ",")

	return messageBasePo
}

func MessageBasePoToDto(messageBasePo *MessageBasePo) (*MessageBaseDto) {
	messageBaseDto := &MessageBaseDto{}
	messageBaseDto.Content = messageBasePo.Content
	messageBaseDto.Tos = strings.Split(messageBasePo.Tos, ",")
	messageBaseDto.ToType = messageBasePo.ToType
	messageBaseDto.ScheduleTime = messageBasePo.ScheduleTime.Unix()
	messageBaseDto.ExpireTime = messageBasePo.ExpireTime.Unix()
	messageBaseDto.Extra = messageBasePo.Extra
	messageBaseDto.SystemAlias = messageBasePo.SystemAlias
	messageBaseDto.FinishedCallbackUrls = strings.Split(messageBasePo.FinishedCallbackUrls, ",")
	messageBaseDto.Id = messageBasePo.Id
	messageBaseDto.SendChannel = messageBasePo.SendChannel
	messageBaseDto.Finished = messageBasePo.Finished
	messageBaseDto.FinishedTime = messageBasePo.FinishedTime
	messageBaseDto.States = messageBasePo.States
	messageBaseDto.ErrorMessages = messageBasePo.ErrorMessages
	messageBaseDto.CreatedTime = messageBasePo.CreatedTime.Unix()
	messageBaseDto.UpdatedTime = messageBasePo.UpdatedTime.Unix()

	return messageBaseDto
}
