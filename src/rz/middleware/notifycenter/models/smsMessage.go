package models

import "strings"

type SmsMessageDto struct {
	MessageBaseDto

	Parameters         []string `json:"parameters"`
	Sign               string   `json:"sign"`
	TemplateId         int      `json:"templateId"`
	NationCode         string   `json:"nationCode"`
	IdentifyingCode    string   `json:"identifyingCode"`
	ExpireCallbackUrls []string `json:"expireCallbackUrls"`
}

type SmsMessagePo struct {
	MessageBasePo

	Parameters         string `gorm:"column:parameters"`
	Sign               string `gorm:"column:sign"`
	TemplateId         int    `gorm:"column:templateId"`
	NationCode         string `gorm:"column:nationCode"`
	IdentifyingCode    string `gorm:"column:identifyingCode"`
	ExpireCallbackUrls string `gorm:"column:expireCallbackUrls"`
}

func (*SmsMessagePo) TableName() string {
	return "smsMessagePo"
}

func SmsMessageDtoToPo(smsMessageDto *SmsMessageDto) (*SmsMessagePo) {
	smsMessagePo := &SmsMessagePo{}
	smsMessagePo.MessageBasePo = *MessageBaseDtoToPo(&smsMessageDto.MessageBaseDto)
	smsMessagePo.Parameters = strings.Join(smsMessageDto.Parameters, ",")
	smsMessagePo.Sign = smsMessageDto.Sign
	smsMessagePo.TemplateId = smsMessageDto.TemplateId
	smsMessagePo.NationCode = smsMessageDto.NationCode
	smsMessagePo.IdentifyingCode = smsMessageDto.IdentifyingCode
	smsMessagePo.ExpireCallbackUrls = strings.Join(smsMessageDto.ExpireCallbackUrls, ",")

	return smsMessagePo
}

func SmsMessagePoToDto(smsMessagePo *SmsMessagePo) (*SmsMessageDto) {
	smsMessageDto := &SmsMessageDto{}
	smsMessageDto.MessageBaseDto = *MessageBasePoToDto(&smsMessagePo.MessageBasePo)
	smsMessageDto.Parameters = strings.Split(smsMessagePo.Parameters, ",")
	smsMessageDto.Sign = smsMessagePo.Sign
	smsMessageDto.TemplateId = smsMessagePo.TemplateId
	smsMessageDto.NationCode = smsMessagePo.NationCode
	smsMessageDto.NationCode = smsMessagePo.NationCode
	smsMessageDto.IdentifyingCode = smsMessagePo.IdentifyingCode
	smsMessageDto.ExpireCallbackUrls = strings.Split(smsMessagePo.ExpireCallbackUrls, ",")

	return smsMessageDto
}
