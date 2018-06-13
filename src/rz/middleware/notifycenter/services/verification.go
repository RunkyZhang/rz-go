package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/exceptions"
)

func verifyMessageDto(messageBaseDto *models.MessageBaseDto) (error) {
	if 0 == len(messageBaseDto.Tos) {
		return exceptions.ErrorTosEmpty()
	}

	return nil
}

func VerifyMailMessageDto(mailMessageDto *models.MailMessageDto) (error) {
	if nil == mailMessageDto {
		return exceptions.DtoNull()
	}

	err := verifyMessageDto(&mailMessageDto.MessageBaseDto)
	if nil != err {
		return err
	}

	if common.IsStringBlank(mailMessageDto.Subject) {
		return exceptions.SubjectBlank()
	}

	return nil
}

func VerifySmsMessageDto(smsMessageDto *models.SmsMessageDto) (error) {
	if nil == smsMessageDto {
		return exceptions.DtoNull()
	}

	err := verifyMessageDto(&smsMessageDto.MessageBaseDto)
	if nil != err {
		return err
	}

	return nil
}
