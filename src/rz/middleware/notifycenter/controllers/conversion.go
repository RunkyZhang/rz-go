package controllers

import (
	"io"
	"rz/middleware/notifycenter/models"
	"encoding/json"
	"rz/middleware/notifycenter/exceptions"
)

func convertToMailMessageDto(body io.ReadCloser) (interface{}, error) {
	var mailMessageDto models.MailMessageDto

	decoder := json.NewDecoder(body)
	exception := decoder.Decode(&mailMessageDto)
	if nil != exception {
		return nil, exceptions.InvalidDtoType
	} else {
		return mailMessageDto, nil
	}
}

func convertToSmsMessageDto(body io.ReadCloser) (interface{}, error) {
	var smsMessageDto models.SmsMessageDto

	decoder := json.NewDecoder(body)
	exception := decoder.Decode(&smsMessageDto)
	if nil != exception {
		return nil, exceptions.InvalidDtoType
	} else {
		return smsMessageDto, nil
	}
}
