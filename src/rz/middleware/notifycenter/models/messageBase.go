package models

import (
	"time"
	"strings"

	"rz/middleware/notifycenter/enumerations"
)

type MessageBaseDto struct {
	Id           int64                    `json:"id,string"`
	Content      string                   `json:"content"`
	Tos          []string                 `json:"tos"`
	ToType       enumerations.ToType      `json:"toType"`
	ScheduleTime int64                    `json:"scheduleTime,string"`
	Extra        string                   `json:"extra"`
	SystemAlias  string                   `json:"systemAlias"`
	SendChannel  enumerations.SendChannel `json:"sendChannel"`
	ContentType  int                      `json:"contentType"`

	FinishedCallbackUrls []string `json:"finishedCallbackUrls"`
	ExpireCallbackUrls   []string `json:"expireCallbackUrls"`
	ExpireTime           int64    `json:"expireTime,string"`
	Disable              bool     `json:"disable"`
	Finished             bool     `json:"finished"`
	FinishedTime         int64    `json:"finishedTime,string"`
	States               string   `json:"states"`
	ErrorMessages        string   `json:"errorMessages"`
	ProviderId           string   `json:"providerId"`
	CreatedTime          int64    `json:"createdTime,string"`
	UpdatedTime          int64    `json:"updatedTime,string"`
}

type MessageBasePo struct {
	PoBase
	CallbackBasePo

	Id           int64                    `gorm:"column:id;primary_key"`
	Content      string                   `gorm:"column:content"`
	Tos          string                   `gorm:"column:tos"`
	ToType       enumerations.ToType      `gorm:"column:toType"`
	ScheduleTime time.Time                `gorm:"column:scheduleTime"`
	Extra        string                   `gorm:"column:extra"`
	SystemAlias  string                   `gorm:"column:systemAlias"`
	SendChannel  enumerations.SendChannel `gorm:"column:sendChannel"`
	ContentType  int                      `gorm:"column:contentType"`
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
	if "" != messageBasePo.Tos {
		messageBaseDto.Tos = strings.Split(messageBasePo.Tos, ",")
	}
	messageBaseDto.ToType = messageBasePo.ToType
	messageBaseDto.ScheduleTime = messageBasePo.ScheduleTime.Unix()
	messageBaseDto.ExpireTime = messageBasePo.ExpireTime.Unix()
	messageBaseDto.Extra = messageBasePo.Extra
	messageBaseDto.SystemAlias = messageBasePo.SystemAlias
	if "" != messageBasePo.FinishedCallbackUrls {
		messageBaseDto.FinishedCallbackUrls = strings.Split(messageBasePo.FinishedCallbackUrls, ",")
	}
	messageBaseDto.Id = messageBasePo.Id
	messageBaseDto.SendChannel = messageBasePo.SendChannel
	messageBaseDto.ProviderId = messageBasePo.ProviderId
	messageBaseDto.Disable = messageBasePo.Disable
	messageBaseDto.Finished = messageBasePo.Finished
	messageBaseDto.FinishedTime = messageBasePo.FinishedTime.Unix()
	messageBaseDto.States = messageBasePo.States
	messageBaseDto.ErrorMessages = messageBasePo.ErrorMessages
	messageBaseDto.CreatedTime = messageBasePo.CreatedTime.Unix()
	messageBaseDto.UpdatedTime = messageBasePo.UpdatedTime.Unix()

	return messageBaseDto
}
