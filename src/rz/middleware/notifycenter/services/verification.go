package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/exceptions"
)

func verifyMessageDto(messageDto *models.MessageDto) (error) {
	if 0 == len(messageDto.Tos) {
		return exceptions.ErrorTosEmpty
	}

	return nil
}

func VerifyMailMessageDto(mailMessageDto *models.MailMessageDto) (error) {
	if nil == mailMessageDto {
		return exceptions.DtoNull
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

func VerifySmsMessageDto(smsMessageDto *models.SmsMessageDto) (error) {
	if nil == smsMessageDto {
		return exceptions.DtoNull
	}

	err := verifyMessageDto(&smsMessageDto.MessageDto)
	if nil != err {
		return err
	}

	return nil
}
