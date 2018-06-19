package controllers

import (
	"io"
	"encoding/json"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
)

func ConvertToMailMessageDto(body io.ReadCloser) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var mailMessageDto models.MailMessageDto
	return convertToDto(body, &mailMessageDto)
}

func ConvertToSmsMessageDto(body io.ReadCloser) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var smsMessageDto models.SmsMessageDto
	return convertToDto(body, &smsMessageDto)
}

func ConvertToSmsCallbackMessageDto(body io.ReadCloser) (interface{}, error) {
	err := common.Assert.IsNotNilToError(body, "body")
	if nil != err {
		return nil, err
	}

	var smsUserCallbackMessageRequestExternalDto external.SmsUserCallbackMessageRequestExternalDto
	return convertToDto(body, &smsUserCallbackMessageRequestExternalDto)
}

func convertToDto(body io.ReadCloser, messageDto interface{}) (interface{}, error) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&messageDto)
	if nil != err {
		return nil, exceptions.InvalidDtoType().AttachError(err)
	}

	return messageDto, nil
}
