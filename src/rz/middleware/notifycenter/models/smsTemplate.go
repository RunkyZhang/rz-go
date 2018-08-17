package models

import (
	"strings"
	"rz/middleware/notifycenter/enumerations"
)

type SmsTemplateDto struct {
	Id                int                          `json:"id"`
	TencentTemplateId int                          `json:"tencentTemplateId"`
	Context           string                       `json:"context"`
	Extend            int                          `json:"extend"`
	UserCallbackUrls  []string                     `json:"userCallbackUrls"`
	Pattern           string                       `json:"pattern"`
	Type              enumerations.SmsTemplateType `json:"type"`
	Sign              string                       `json:"sign"`
	DahanSignCode     int                          `json:"dahanSignCode"`
	TencentContext    string                       `json:"tencentContext"`
	DahanContext      string                       `json:"dahanContext"`
}

type SmsTemplatePo struct {
	PoBase

	Id                int                          `gorm:"column:id;primary_key;auto_increment"`
	TencentTemplateId int                          `gorm:"column:tencentTemplateId"`
	Context           string                       `gorm:"column:context"`
	Extend            int                          `gorm:"column:extend"`
	UserCallbackUrls  string                       `gorm:"column:userCallbackUrls"`
	Pattern           string                       `gorm:"column:pattern"`
	Type              enumerations.SmsTemplateType `gorm:"column:type"`
	Sign              string                       `gorm:"column:sign"`
	DahanSignCode     int                          `gorm:"column:dahanSignCode"`
	TencentContext    string                       `gorm:"column:tencentContext"`
	DahanContext      string                       `gorm:"column:dahanContext"`
}

func SmsTemplateDtoToPo(smsTemplateDto *SmsTemplateDto) (*SmsTemplatePo) {
	if nil == smsTemplateDto {
		return nil
	}

	smsTemplatePo := &SmsTemplatePo{}
	smsTemplatePo.Id = smsTemplateDto.Id
	smsTemplatePo.TencentTemplateId = smsTemplateDto.TencentTemplateId
	smsTemplatePo.Context = smsTemplateDto.Context
	smsTemplatePo.Extend = smsTemplateDto.Extend
	smsTemplatePo.Pattern = smsTemplateDto.Pattern
	smsTemplatePo.UserCallbackUrls = strings.Join(smsTemplateDto.UserCallbackUrls, ",")
	smsTemplatePo.Type = smsTemplateDto.Type
	smsTemplatePo.Sign = smsTemplateDto.Sign
	smsTemplatePo.DahanSignCode = smsTemplateDto.DahanSignCode
	smsTemplatePo.TencentContext = smsTemplateDto.TencentContext
	smsTemplatePo.DahanContext = smsTemplateDto.DahanContext

	return smsTemplatePo
}

func SmsTemplatePoToDto(smsTemplatePo *SmsTemplatePo) (*SmsTemplateDto) {
	if nil == smsTemplatePo {
		return nil
	}

	smsTemplateDto := &SmsTemplateDto{}
	smsTemplateDto.Id = smsTemplatePo.Id
	smsTemplateDto.TencentTemplateId = smsTemplatePo.TencentTemplateId
	smsTemplateDto.Context = smsTemplatePo.Context
	smsTemplateDto.Extend = smsTemplatePo.Extend
	smsTemplateDto.Pattern = smsTemplatePo.Pattern
	if "" != smsTemplatePo.UserCallbackUrls {
		smsTemplateDto.UserCallbackUrls = strings.Split(smsTemplatePo.UserCallbackUrls, ",")
	}
	smsTemplateDto.Type = smsTemplatePo.Type
	smsTemplateDto.Sign = smsTemplatePo.Sign
	smsTemplateDto.DahanSignCode = smsTemplatePo.DahanSignCode
	smsTemplateDto.TencentContext = smsTemplatePo.TencentContext
	smsTemplateDto.DahanContext = smsTemplatePo.DahanContext

	return smsTemplateDto
}

func SmsTemplatePosToDtos(smsTemplatePos []*SmsTemplatePo) ([]*SmsTemplateDto) {
	if nil == smsTemplatePos {
		return nil
	}

	var smsTemplateDtos []*SmsTemplateDto
	for _, smsMessagePo := range smsTemplatePos {
		smsTemplateDtos = append(smsTemplateDtos, SmsTemplatePoToDto(smsMessagePo))
	}

	return smsTemplateDtos
}
