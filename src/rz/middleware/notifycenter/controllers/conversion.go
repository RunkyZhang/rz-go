package controllers

import (
	"encoding/json"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
)

func ConvertToMailMessageDto(body []byte) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var dto models.MailMessageDto
	return convertToDto(body, &dto)
}

func ConvertToSmsMessageDto(body []byte) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var dto models.SmsMessageDto
	return convertToDto(body, &dto)
}

func ConvertToQueryMessagesByIdsRequestDto(body []byte) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var dto models.QueryMessagesByIdsRequestDto
	return convertToDto(body, &dto)
}

func ConvertToDisableMessageRequestDto(body []byte) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var dto models.DisableMessageRequestDto
	return convertToDto(body, &dto)
}

func ConvertToSmsCallbackMessageDto(body []byte) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var dto external.SmsUserCallbackMessageRequestExternalDto
	return convertToDto(body, &dto)
}

func ConvertToSmsTemplateDto(body []byte) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var dto models.SmsTemplateDto
	return convertToDto(body, &dto)
}

func ConvertToQuerySmsUserMessagesRequestDto(body []byte) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var dto models.QuerySmsUserMessagesRequestDto
	return convertToDto(body, &dto)
}

func convertToDto(body []byte, messageDto interface{}) (interface{}, error) {
	err := json.Unmarshal(body, &messageDto)
	if nil != err {
		return nil, exceptions.InvalidDtoType().AttachError(err)
	}

	return messageDto, nil
}
