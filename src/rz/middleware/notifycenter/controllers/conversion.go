package controllers

import (
	"io"
	"encoding/json"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/models/external"
)

func ConvertToMailMessageDto(body io.ReadCloser) (interface{}, error) {
	var mailMessageDto models.MailMessageDto

	return convertToDto(body, &mailMessageDto)
}

func ConvertToSmsMessageDto(body io.ReadCloser) (interface{}, error) {
	var smsMessageDto models.SmsMessageDto

	return convertToDto(body, &smsMessageDto)
}

func ConvertToSmsCallbackMessageDto(body io.ReadCloser) (interface{}, error) {
	var smsUserCallbackMessageRequestExternalDto external.SmsUserCallbackMessageRequestExternalDto

	return convertToDto(body, &smsUserCallbackMessageRequestExternalDto)
}

func convertToDto(body io.ReadCloser, messageDto interface{}) (interface{}, error) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&messageDto)
	if nil == err {
		return messageDto, nil
	}

	return nil, exceptions.InvalidDtoType
}
