package models

type SystemAliasPermissionDto struct {
	SystemAlias        string `json:"systemAlias"`
	SmsPermission      int    `json:"smsPermission"`
	MailPermission     int    `json:"mailPermission"`
	SmsDayFrequency    int    `json:"smsDayFrequency"`
	SmsHourFrequency   int    `json:"smsHourFrequency"`
	SmsMinuteFrequency int    `json:"smsMinuteFrequency"`
	CreatedTime        int64  `json:"createdTime,string"`
	UpdatedTime        int64  `json:"updatedTime,string"`
}

type SystemAliasPermissionPo struct {
	PoBase

	SystemAlias        string `gorm:"column:id;primary_key"`
	SmsPermission      int    `gorm:"column:smsPermission"`
	MailPermission     int    `gorm:"column:mailPermission"`
	SmsDayFrequency    int    `gorm:"column:smsDayFrequency"`
	SmsHourFrequency   int    `gorm:"column:smsHourFrequency"`
	SmsMinuteFrequency int    `gorm:"column:smsMinuteFrequency"`
}

func SystemAliasPermissionDtoToPo(systemAliasPermissionDto *SystemAliasPermissionDto) (*SystemAliasPermissionPo) {
	if nil == systemAliasPermissionDto {
		return nil
	}

	systemAliasPermissionPo := &SystemAliasPermissionPo{}
	systemAliasPermissionPo.SystemAlias = systemAliasPermissionDto.SystemAlias
	systemAliasPermissionPo.SmsPermission = systemAliasPermissionDto.SmsPermission
	systemAliasPermissionPo.MailPermission = systemAliasPermissionDto.MailPermission
	systemAliasPermissionPo.SmsDayFrequency = systemAliasPermissionDto.SmsDayFrequency
	systemAliasPermissionPo.SmsHourFrequency = systemAliasPermissionDto.SmsHourFrequency
	systemAliasPermissionPo.SmsMinuteFrequency = systemAliasPermissionDto.SmsMinuteFrequency

	return systemAliasPermissionPo
}

func SystemAliasPermissionPoToDto(systemAliasPermissionPo *SystemAliasPermissionPo) (*SystemAliasPermissionDto) {
	if nil == systemAliasPermissionPo {
		return nil
	}

	systemAliasPermissionDto := &SystemAliasPermissionDto{}
	systemAliasPermissionDto.SystemAlias = systemAliasPermissionPo.SystemAlias
	systemAliasPermissionDto.SmsPermission = systemAliasPermissionPo.SmsPermission
	systemAliasPermissionDto.MailPermission = systemAliasPermissionPo.MailPermission
	systemAliasPermissionDto.SmsDayFrequency = systemAliasPermissionPo.SmsDayFrequency
	systemAliasPermissionDto.SmsHourFrequency = systemAliasPermissionPo.SmsHourFrequency
	systemAliasPermissionDto.SmsMinuteFrequency = systemAliasPermissionPo.SmsMinuteFrequency
	systemAliasPermissionDto.CreatedTime = systemAliasPermissionPo.CreatedTime.Unix()
	systemAliasPermissionDto.UpdatedTime = systemAliasPermissionPo.UpdatedTime.Unix()

	return systemAliasPermissionDto
}

func SystemAliasPermissionPosToDtos(systemAliasPermissionPos []*SystemAliasPermissionPo) ([]*SystemAliasPermissionDto) {
	if nil == systemAliasPermissionPos {
		return nil
	}

	var systemAliasPermissionDtos []*SystemAliasPermissionDto
	for _, systemAliasPermissionPo := range systemAliasPermissionPos {
		systemAliasPermissionDtos = append(systemAliasPermissionDtos, SystemAliasPermissionPoToDto(systemAliasPermissionPo))
	}

	return systemAliasPermissionDtos
}
