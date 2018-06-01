package controllers

import (
	"io"
	"encoding/json"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
)

func ConvertToMailMessageDto(body io.ReadCloser) (interface{}, error) {
	var mailMessageDto models.MailMessageDto

	return convertToMessageDto(body, &mailMessageDto)
}

func ConvertToSmsMessageDto(body io.ReadCloser) (interface{}, error) {
	var smsMessageDto models.SmsMessageDto

	return convertToMessageDto(body, &smsMessageDto)
}

func convertToMessageDto(body io.ReadCloser, messageDto interface{}) (interface{}, error) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&messageDto)
	if nil == err {
		return messageDto, nil
	}

	return nil, exceptions.InvalidDtoType
}
