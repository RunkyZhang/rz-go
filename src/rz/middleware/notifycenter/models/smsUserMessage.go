package models

import (
	"time"
)

type SmsUserMessageDto struct {
	Id            int    `json:"id"`
	NationCode    string `json:"nationCode"`
	PhoneNumber   string `json:"phoneNumber"`
	TemplateId    int    `json:"templateId"`
	Extend        int    `json:"extend"`
	Content       string `json:"content"`
	Sign          string `json:"sign"`
	Time          int64  `json:"time"`
	Finished      bool   `json:"finished"`
	FinishedTime  int64  `json:"finishedTime"`
	ErrorMessages string `json:"errorMessages"`
	States        string `json:"states"`
	CreatedTime   int64  `json:"createdTime"`
	UpdatedTime   int64  `json:"updatedTime"`
}

type SmsUserMessagePo struct {
	PoBase
	CallbackBasePo

	NationCode  string `gorm:"column:nationCode"`
	PhoneNumber string `gorm:"column:phoneNumber"`
	TemplateId  int    `gorm:"column:templateId"`
	Extend      int    `gorm:"column:extend"`
	Content     string `gorm:"column:content"`
	Sign        string `gorm:"column:sign"`
	Time        int64  `gorm:"column:time"`
}

func SmsUserMessageDtoToPo(smsUserMessageDto *SmsUserMessageDto) (*SmsUserMessagePo) {
	smsUserMessagePo := &SmsUserMessagePo{}
	smsUserMessagePo.Id = smsUserMessageDto.Id
	smsUserMessagePo.NationCode = smsUserMessageDto.NationCode
	smsUserMessagePo.PhoneNumber = smsUserMessageDto.PhoneNumber
	smsUserMessagePo.TemplateId = smsUserMessageDto.TemplateId
	smsUserMessagePo.Content = smsUserMessageDto.Content
	smsUserMessagePo.Sign = smsUserMessageDto.Sign
	smsUserMessagePo.Time = smsUserMessageDto.Time
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
	smsUserMessageDto.Time = smsUserMessagePo.Time
	smsUserMessageDto.Finished = smsUserMessagePo.Finished
	smsUserMessageDto.FinishedTime = smsUserMessagePo.FinishedTime.Unix()
	smsUserMessageDto.ErrorMessages = smsUserMessagePo.ErrorMessages
	smsUserMessageDto.CreatedTime = smsUserMessagePo.CreatedTime.Unix()
	smsUserMessageDto.UpdatedTime = smsUserMessagePo.UpdatedTime.Unix()

	return smsUserMessageDto
}
