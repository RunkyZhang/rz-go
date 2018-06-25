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

	var mailMessageDto models.MailMessageDto
	return convertToDto(body, &mailMessageDto)
}

func ConvertToSmsMessageDto(body []byte) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var smsMessageDto models.SmsMessageDto
	return convertToDto(body, &smsMessageDto)
}

func ConvertToSmsCallbackMessageDto(body []byte) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var smsUserCallbackMessageRequestExternalDto external.SmsUserCallbackMessageRequestExternalDto
	return convertToDto(body, &smsUserCallbackMessageRequestExternalDto)
}

func convertToDto(body []byte, messageDto interface{}) (interface{}, error) {
	err := json.Unmarshal(body, &messageDto)
	if nil != err {
		return nil, exceptions.InvalidDtoType().AttachError(err)
	}

	return messageDto, nil
}
