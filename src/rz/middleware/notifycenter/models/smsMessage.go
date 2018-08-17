package models

import "strings"

type SmsMessageDto struct {
	MessageBaseDto

	Parameters      []string `json:"parameters"`
	Sign            string   `json:"sign"`
	TemplateId      int      `json:"templateId"`
	NationCode      string   `json:"nationCode"`
	IdentifyingCode string   `json:"identifyingCode"`
}

type SmsMessagePo struct {
	MessageBasePo

	Parameters      string `gorm:"column:parameters"`
	Sign            string `gorm:"column:sign"`
	TemplateId      int    `gorm:"column:templateId"`
	NationCode      string `gorm:"column:nationCode"`
	IdentifyingCode string `gorm:"column:identifyingCode"`
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
	if "" != smsMessagePo.Parameters {
		smsMessageDto.Parameters = strings.Split(smsMessagePo.Parameters, ",")
	}
	smsMessageDto.Sign = smsMessagePo.Sign
	smsMessageDto.TemplateId = smsMessagePo.TemplateId
	smsMessageDto.NationCode = smsMessagePo.NationCode
	smsMessageDto.NationCode = smsMessagePo.NationCode
	smsMessageDto.IdentifyingCode = smsMessagePo.IdentifyingCode
	if "" != smsMessagePo.ExpireCallbackUrls {
		smsMessageDto.ExpireCallbackUrls = strings.Split(smsMessagePo.ExpireCallbackUrls, ",")
	}

	return smsMessageDto
}

func SmsMessagePosToDtos(smsMessagePos []*SmsMessagePo) ([]*SmsMessageDto) {
	if nil == smsMessagePos {
		return nil
	}

	var smsMessageDtos []*SmsMessageDto
	for _, smsMessagePo := range smsMessagePos {
		smsMessageDtos = append(smsMessageDtos, SmsMessagePoToDto(smsMessagePo))
	}

	return smsMessageDtos
}
