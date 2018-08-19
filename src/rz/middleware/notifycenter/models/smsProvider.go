package models

type SmsProviderDto struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	ContentTypes       int    `json:"contentTypes"`
	PassportJson       string `json:"passportJson"`
	Priority           int    `json:"priority"`
	Weighted           string `json:"weighted"`
	Description        string `json:"description"`
	Url1               string `json:"url1"`
	Url2               string `json:"url2"`
	CallbackUrl        string `json:"callbackUrl"`
	ManagementUrl      string `json:"managementUrl"`
	ManagementUser     string `json:"managementUser"`
	ManagementPassword string `json:"managementPassword"`
	Disable            bool   `json:"disable"`
}

type SmsProviderPo struct {
	PoBase

	Id                 string `gorm:"column:id;primary_key"`
	Name               string `gorm:"column:name"`
	ContentTypes       int    `gorm:"column:contentTypes"`
	PassportJson       string `gorm:"column:passportJson"`
	Priority           int    `gorm:"column:priority"`
	Weighted           string `gorm:"column:weighted"`
	Description        string `gorm:"column:description"`
	Url1               string `gorm:"column:url1"`
	Url2               string `gorm:"column:url2"`
	CallbackUrl        string `gorm:"column:callbackUrl"`
	ManagementUrl      string `gorm:"column:managementUrl"`
	ManagementUser     string `gorm:"column:managementUser"`
	ManagementPassword string `gorm:"column:managementPassword"`
	Disable            bool   `gorm:"column:disable"`
}

func SmsProviderDtoToPo(smsProviderDto *SmsProviderDto) (*SmsProviderPo) {
	if nil == smsProviderDto {
		return nil
	}

	smsProviderPo := &SmsProviderPo{}
	smsProviderPo.Id = smsProviderDto.Id
	smsProviderPo.Name = smsProviderDto.Name
	smsProviderPo.ContentTypes = smsProviderDto.ContentTypes
	smsProviderPo.PassportJson = smsProviderDto.PassportJson
	smsProviderPo.Priority = smsProviderDto.Priority
	smsProviderPo.Weighted = smsProviderDto.Weighted
	smsProviderPo.Description = smsProviderDto.Description
	smsProviderPo.Url1 = smsProviderDto.Url1
	smsProviderPo.Url2 = smsProviderDto.Url2
	smsProviderPo.CallbackUrl = smsProviderDto.CallbackUrl
	smsProviderPo.ManagementUrl = smsProviderDto.ManagementUrl
	smsProviderPo.ManagementUser = smsProviderDto.ManagementUser
	smsProviderPo.ManagementPassword = smsProviderDto.ManagementPassword
	smsProviderPo.Disable = smsProviderDto.Disable

	return smsProviderPo
}

func SmsProviderPoToDto(smsProviderPo *SmsProviderPo) (*SmsProviderDto) {
	if nil == smsProviderPo {
		return nil
	}

	smsProviderDto := &SmsProviderDto{}
	smsProviderDto.Id = smsProviderPo.Id
	smsProviderDto.Name = smsProviderPo.Name
	smsProviderDto.ContentTypes = smsProviderPo.ContentTypes
	smsProviderDto.PassportJson = smsProviderPo.PassportJson
	smsProviderDto.Priority = smsProviderPo.Priority
	smsProviderDto.Weighted = smsProviderPo.Weighted
	smsProviderDto.Description = smsProviderPo.Description
	smsProviderDto.Url1 = smsProviderPo.Url1
	smsProviderDto.Url2 = smsProviderPo.Url2
	smsProviderDto.CallbackUrl = smsProviderPo.CallbackUrl
	smsProviderDto.ManagementUrl = smsProviderPo.ManagementUrl
	smsProviderDto.ManagementUser = smsProviderPo.ManagementUser
	smsProviderDto.ManagementPassword = smsProviderPo.ManagementPassword
	smsProviderDto.Disable = smsProviderPo.Disable

	return smsProviderDto
}

func SmsProviderPosToDtos(smsProviderPos []*SmsProviderPo) ([]*SmsProviderDto) {
	if nil == smsProviderPos {
		return nil
	}

	var smsProviderDtos []*SmsProviderDto
	for _, smsProviderPo := range smsProviderPos {
		smsProviderDtos = append(smsProviderDtos, SmsProviderPoToDto(smsProviderPo))
	}

	return smsProviderDtos
}
