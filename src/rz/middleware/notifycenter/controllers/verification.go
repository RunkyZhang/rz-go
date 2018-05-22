package controllers

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
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

	mailMessageDto, ok := dto.(*models.MailMessageDto)
	if !ok {
		return exceptions.InvalidDtoType
	}

	err := verifyMessageDto(&mailMessageDto.MessageDto)
	if nil != err {
		return err
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

	err := verifyMessageDto(&smsMessageDto.MessageDto)
	if nil != err {
		return err
	}

	return nil
}
