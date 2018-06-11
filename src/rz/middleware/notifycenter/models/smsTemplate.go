package models

import "strings"

type SmsTemplateDto struct {
	Id               int      `json:"id"`
	Extend           int      `json:"extend"`
	UserCallbackUrls []string `json:"userCallbackUrls"`
	Pattern          string   `json:"pattern"`
}

type SmsTemplatePo struct {
	PoBase

	Id               int    `json:"id"`
	Extend           int    `json:"extend"`
	UserCallbackUrls string `json:"userCallbackUrls"`
	Pattern          string `json:"pattern"`
}

func SmsTemplateDtoToPo(smsTemplateDto *SmsTemplateDto) (*SmsTemplatePo) {
	smsTemplatePo := &SmsTemplatePo{}
	smsTemplatePo.Id = smsTemplateDto.Id
	smsTemplatePo.Extend = smsTemplateDto.Extend
	smsTemplatePo.Pattern = smsTemplateDto.Pattern
	smsTemplatePo.UserCallbackUrls = strings.Join(smsTemplateDto.UserCallbackUrls, ",")

	return smsTemplatePo
}

func SmsTemplatePoToDto(smsTemplatePo *SmsTemplatePo) (*SmsTemplateDto) {
	smsTemplateDto := &SmsTemplateDto{}
	smsTemplateDto.Id = smsTemplatePo.Id
	smsTemplateDto.Extend = smsTemplatePo.Extend
	smsTemplateDto.Pattern = smsTemplatePo.Pattern
	smsTemplateDto.UserCallbackUrls = strings.Split(smsTemplatePo.UserCallbackUrls, ",")

	return smsTemplateDto
}
