package controllers

import (
	"io"
	"rz/middleware/notifycenter/models"
	"encoding/json"
	"rz/middleware/notifycenter/exceptions"
)

func convertToMailMessageDto(body io.ReadCloser) (interface{}, error) {
	var mailMessageDto models.MailMessageDto

	return convertToMessageDto(body, &mailMessageDto)
}

func convertToSmsMessageDto(body io.ReadCloser) (interface{}, error) {
	var smsMessageDto models.SmsMessageDto

	return convertToMessageDto(body, &smsMessageDto)
}

func convertToMessageDto(body io.ReadCloser, messageDto interface{}) (interface{}, error) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&messageDto)
	if nil != err {
		return nil, exceptions.InvalidDtoType
	} else {
		return messageDto, nil
	}
}
