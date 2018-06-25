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

func SmsMessageDtoToPo(smsMessageDto *SmsMessageDto) (*SmsMessagePo) {
	if nil == smsMessageDto {
		return nil
	}

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
	if nil == smsMessagePo {
		return nil
	}

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

type QuerySmsMessageRequestDto struct {
	Id          int    `json:"id"`
	SystemAlias string `json:"systemAlias"`
	CreatedTime int64  `json:"createdTime"`
}
