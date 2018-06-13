package models

import (
	"strings"
	"rz/middleware/notifycenter/enumerations"
)

type SmsTemplateDto struct {
	Id               int                          `json:"id"`
	Extend           int                          `json:"extend"`
	UserCallbackUrls []string                     `json:"userCallbackUrls"`
	Pattern          string                       `json:"pattern"`
	Type             enumerations.SmsTemplateType `json:"type"`
}

type SmsTemplatePo struct {
	PoBase

	Id               int                          `gorm:"column:id;primary_key"`
	Extend           int                          `gorm:"column:extend"`
	UserCallbackUrls string                       `gorm:"column:userCallbackUrls"`
	Pattern          string                       `gorm:"column:pattern"`
	Type             enumerations.SmsTemplateType `gorm:"column:type"`
}

func SmsTemplateDtoToPo(smsTemplateDto *SmsTemplateDto) (*SmsTemplatePo) {
	if nil == smsTemplateDto {
		return nil
	}

	smsTemplatePo := &SmsTemplatePo{}
	smsTemplatePo.Id = smsTemplateDto.Id
	smsTemplatePo.Extend = smsTemplateDto.Extend
	smsTemplatePo.Pattern = smsTemplateDto.Pattern
	smsTemplatePo.UserCallbackUrls = strings.Join(smsTemplateDto.UserCallbackUrls, ",")
	smsTemplatePo.Type = smsTemplateDto.Type

	return smsTemplatePo
}

func SmsTemplatePoToDto(smsTemplatePo *SmsTemplatePo) (*SmsTemplateDto) {
	if nil == smsTemplatePo {
		return nil
	}

	smsTemplateDto := &SmsTemplateDto{}
	smsTemplateDto.Id = smsTemplatePo.Id
	smsTemplateDto.Extend = smsTemplatePo.Extend
	smsTemplateDto.Pattern = smsTemplatePo.Pattern
	smsTemplateDto.UserCallbackUrls = strings.Split(smsTemplatePo.UserCallbackUrls, ",")
	smsTemplateDto.Type = smsTemplatePo.Type

	return smsTemplateDto
}
