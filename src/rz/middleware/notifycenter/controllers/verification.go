package controllers

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
	"fmt"
)

func verifyMessageDto(messageDto *models.MessageDto) (error) {
	if 0 == len(messageDto.Tos) {
		return exceptions.ErrorTosEmpty
	}

	return nil
}

func verifyMailMessageDto(dto interface{}) (error) {
	if nil == dto {
		return exceptions.DtoNull
	}

	fmt.Println(&dto)
	mailMessageDto, ok := dto.(*models.MailMessageDto)
	if !ok {
		return exceptions.InvalidDtoType
	}

	exception := verifyMessageDto(&mailMessageDto.MessageDto)
	if nil != exception {
		return exception
	}

	if common.IsStringBlank(mailMessageDto.Subject) {
		return exceptions.SubjectBlank
	}

	return nil
}

func verifySmsMessageDto(dto interface{}) (error) {
	if nil == dto {
		return exceptions.DtoNull
	}

	smsMessageDto, ok := dto.(*models.SmsMessageDto)
	if !ok {
		return exceptions.InvalidDtoType
	}

	exception := verifyMessageDto(&smsMessageDto.MessageDto)
	if nil != exception {
		return exception
	}

	return nil
}
