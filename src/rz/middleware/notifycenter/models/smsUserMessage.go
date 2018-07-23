package models

import (
	"time"
)

type SmsUserMessageDto struct {
	Id           int64  `json:"id,string"`
	NationCode   string `json:"nationCode"`
	PhoneNumber  string `json:"phoneNumber"`
	TemplateId   int    `json:"templateId"`
	Extend       int    `json:"extend"`
	Content      string `json:"content"`
	Sign         string `json:"sign"`
	Time         int64  `json:"time,string"`
	SmsMessageId int64  `json:"time,string"`

	Disable       bool   `json:"disable"`
	Finished      bool   `json:"finished"`
	FinishedTime  int64  `json:"finishedTime,string"`
	ErrorMessages string `json:"errorMessages"`
	States        string `json:"states"`
	CreatedTime   int64  `json:"createdTime,string"`
	UpdatedTime   int64  `json:"updatedTime,string"`
}

type SmsUserMessagePo struct {
	PoBase
	CallbackBasePo

	Id           int64  `gorm:"column:id;primary_key"`
	NationCode   string `gorm:"column:nationCode"`
	PhoneNumber  string `gorm:"column:phoneNumber"`
	TemplateId   int    `gorm:"column:templateId"`
	Extend       int    `gorm:"column:extend"`
	Content      string `gorm:"column:content"`
	Sign         string `gorm:"column:sign"`
	Time         int64  `gorm:"column:time"`
	SmsMessageId int64  `gorm:"column:smsMessageId"`
}

func SmsUserMessageDtoToPo(smsUserMessageDto *SmsUserMessageDto) (*SmsUserMessagePo) {
	if nil == smsUserMessageDto {
		return nil
	}

	smsUserMessagePo := &SmsUserMessagePo{}
	smsUserMessagePo.Id = smsUserMessageDto.Id
	smsUserMessagePo.NationCode = smsUserMessageDto.NationCode
	smsUserMessagePo.PhoneNumber = smsUserMessageDto.PhoneNumber
	smsUserMessagePo.TemplateId = smsUserMessageDto.TemplateId
	smsUserMessagePo.Content = smsUserMessageDto.Content
	smsUserMessagePo.Sign = smsUserMessageDto.Sign
	smsUserMessagePo.Time = smsUserMessageDto.Time
	smsUserMessagePo.SmsMessageId = smsUserMessageDto.SmsMessageId
	smsUserMessagePo.Finished = smsUserMessageDto.Finished
	smsUserMessagePo.FinishedTime = time.Unix(smsUserMessageDto.FinishedTime, 0)
	smsUserMessagePo.ErrorMessages = smsUserMessageDto.ErrorMessages

	return smsUserMessagePo
}

func SmsUserMessagePoToDto(smsUserMessagePo *SmsUserMessagePo) (*SmsUserMessageDto) {
	if nil == smsUserMessagePo {
		return nil
	}

	smsUserMessageDto := &SmsUserMessageDto{}
	smsUserMessageDto.Id = smsUserMessagePo.Id
	smsUserMessageDto.NationCode = smsUserMessagePo.NationCode
	smsUserMessageDto.PhoneNumber = smsUserMessagePo.PhoneNumber
	smsUserMessageDto.TemplateId = smsUserMessagePo.TemplateId
	smsUserMessageDto.Content = smsUserMessagePo.Content
	smsUserMessageDto.Sign = smsUserMessagePo.Sign
	smsUserMessageDto.Time = smsUserMessagePo.Time
	smsUserMessageDto.SmsMessageId = smsUserMessagePo.SmsMessageId
	smsUserMessageDto.Disable = smsUserMessagePo.Disable
	smsUserMessageDto.Finished = smsUserMessagePo.Finished
	smsUserMessageDto.FinishedTime = smsUserMessagePo.FinishedTime.Unix()
	smsUserMessageDto.ErrorMessages = smsUserMessagePo.ErrorMessages
	smsUserMessageDto.States = smsUserMessagePo.States
	smsUserMessageDto.CreatedTime = smsUserMessagePo.CreatedTime.Unix()
	smsUserMessageDto.UpdatedTime = smsUserMessagePo.UpdatedTime.Unix()

	return smsUserMessageDto
}

func SmsUserMessagePosToDtos(smsUserMessagePos []*SmsUserMessagePo) ([]*SmsUserMessageDto) {
	if nil == smsUserMessagePos {
		return nil
	}

	var smsUserMessageDtos []*SmsUserMessageDto
	for _, smsUserMessagePo := range smsUserMessagePos {
		smsUserMessageDtos = append(smsUserMessageDtos, SmsUserMessagePoToDto(smsUserMessagePo))
	}

	return smsUserMessageDtos
}
