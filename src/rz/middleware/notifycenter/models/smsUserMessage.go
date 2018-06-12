package models

import (
	"time"
)

type SmsUserMessageDto struct {
	Id            string `json:"id"`
	NationCode    string `json:"nationCode"`
	PhoneNumber   string `json:"phoneNumber"`
	TemplateId    int    `json:"templateId"`
	Content       string `json:"content"`
	Sign          string `json:"sign"`
	Time          int64  `json:"time"`
	Finished      bool   `json:"finished"`
	FinishedTime  int64  `json:"finishedTime"`
	ErrorMessages string `json:"errorMessages"`

	CreatedTime int64 `json:"createdTime"`
	UpdatedTime int64 `json:"updatedTime"`
}

type SmsUserMessagePo struct {
	PoBase

	Id            string    `gorm:"column:id;primary_key"`
	NationCode    string    `gorm:"column:nationCode"`
	PhoneNumber   string    `gorm:"column:phoneNumber"`
	TemplateId    int       `gorm:"column:templateId"`
	Content       string    `gorm:"column:content"`
	Sign          string    `gorm:"column:sign"`
	Time          time.Time `gorm:"column:time"`
	Finished      bool      `gorm:"column:finished"`
	FinishedTime  time.Time `gorm:"column:finishedTime"`
	ErrorMessages string    `gorm:"column:errorMessages"`
}

func SmsUserMessageDtoToPo(smsUserMessageDto *SmsUserMessageDto) (*SmsUserMessagePo) {
	smsUserMessagePo := &SmsUserMessagePo{}
	smsUserMessagePo.Id = smsUserMessageDto.Id
	smsUserMessagePo.NationCode = smsUserMessageDto.NationCode
	smsUserMessagePo.PhoneNumber = smsUserMessageDto.PhoneNumber
	smsUserMessagePo.TemplateId = smsUserMessageDto.TemplateId
	smsUserMessagePo.Content = smsUserMessageDto.Content
	smsUserMessagePo.Sign = smsUserMessageDto.Sign
	smsUserMessagePo.Time = time.Unix(smsUserMessageDto.Time, 0)
	smsUserMessagePo.Finished = smsUserMessageDto.Finished
	smsUserMessagePo.FinishedTime = time.Unix(smsUserMessageDto.FinishedTime, 0)
	smsUserMessagePo.ErrorMessages = smsUserMessageDto.ErrorMessages

	return smsUserMessagePo
}

func SmsUserMessageToDto(smsUserMessagePo *SmsUserMessagePo) (*SmsUserMessageDto) {
	smsUserMessageDto := &SmsUserMessageDto{}
	smsUserMessageDto.Id = smsUserMessagePo.Id
	smsUserMessageDto.NationCode = smsUserMessagePo.NationCode
	smsUserMessageDto.PhoneNumber = smsUserMessagePo.PhoneNumber
	smsUserMessageDto.TemplateId = smsUserMessagePo.TemplateId
	smsUserMessageDto.Content = smsUserMessagePo.Content
	smsUserMessageDto.Sign = smsUserMessagePo.Sign
	smsUserMessageDto.Time = smsUserMessagePo.Time.Unix()
	smsUserMessageDto.Finished = smsUserMessagePo.Finished
	smsUserMessageDto.FinishedTime = smsUserMessagePo.FinishedTime.Unix()
	smsUserMessageDto.ErrorMessages = smsUserMessagePo.ErrorMessages
	smsUserMessageDto.CreatedTime = smsUserMessagePo.CreatedTime.Unix()
	smsUserMessageDto.UpdatedTime = smsUserMessagePo.UpdatedTime.Unix()

	return smsUserMessageDto
}
