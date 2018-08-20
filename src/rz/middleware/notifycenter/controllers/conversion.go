package controllers

import (
	"encoding/json"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/exceptions"
	"rz/core/common"
)

func ConvertToMailMessageDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto models.MailMessageDto
	return convertToDto(body, &dto)
}

func ConvertToSmsMessageDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto models.SmsMessageDto
	return convertToDto(body, &dto)
}

func ConvertToQueryMessagesByIdsRequestDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto models.QueryMessagesByIdsRequestDto
	return convertToDto(body, &dto)
}

func ConvertToDisableMessageRequestDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto models.DisableMessageRequestDto
	return convertToDto(body, &dto)
}

func ConvertToTencentSmsUserCallbackRequestDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto external.TencentSmsUserCallbackRequestDto
	return convertToDto(body, &dto)
}

func ConvertToDahanSmsUserCallbackRequestDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto external.DahanSmsUserCallbackRequestDto
	return convertToDto(body, &dto)
}

func ConvertToSmsTemplateDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto models.SmsTemplateDto
	return convertToDto(body, &dto)
}

func ConvertToQuerySmsUserMessagesRequestDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto models.QuerySmsUserMessagesRequestDto
	return convertToDto(body, &dto)
}

func ConvertToSystemAliasPermissionDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto models.SystemAliasPermissionDto
	return convertToDto(body, &dto)
}

func ConvertToModifySystemAliasPermissionRequestDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto models.ModifySystemAliasPermissionRequestDto
	return convertToDto(body, &dto)
}

func ConvertToSmsProviderDto(body []byte) (interface{}, error) {
	err := common.Assert.IsTrueToError(nil != body, "nil != body")
	if nil != err {
		return nil, err
	}

	var dto models.SmsProviderDto
	return convertToDto(body, &dto)
}

func convertToDto(body []byte, messageDto interface{}) (interface{}, error) {
	err := json.Unmarshal(body, &messageDto)
	if nil != err {
		return nil, exceptions.InvalidDtoType().AttachError(err)
	}

	return messageDto, nil
}
