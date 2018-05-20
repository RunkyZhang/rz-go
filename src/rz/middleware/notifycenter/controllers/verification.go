package controllers

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
)

func verifyMessageDto(dto interface{}) (error) {
	if nil == dto {
		return exceptions.DtoNull
	}
	messageDto, ok := dto.(*models.MessageDto)
	if !ok {
		return exceptions.InvalidDtoType
	}

	if 0 == len(messageDto.Tos) {
		return exceptions.ErrorTosEmpty
	}

	return nil
}

func verifyMailMessageDto(dto interface{}) (error) {
	exception := verifyMessageDto(dto)
	if nil != exception {
		return exception
	}

	mailMessageDto, ok := dto.(*models.MailMessageDto)
	if !ok {
		return exceptions.InvalidDtoType
	}

	if common.IsStringBlank(mailMessageDto.Subject) {
		return exceptions.SubjectBlank
	}

	return nil
}

func verifySmsMessageDto(dto interface{}) (error) {
	exception := verifyMessageDto(dto)
	if nil != exception {
		return exception
	}

	_, ok := dto.(*models.SmsMessageDto)
	if !ok {
		return exceptions.InvalidDtoType
	}

	return nil
}
