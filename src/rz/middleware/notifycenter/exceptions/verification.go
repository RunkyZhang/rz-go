package exceptions

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
)

func verifyMessageDto(messageDto *models.MessageDto) (error) {
	if 0 == len(messageDto.Tos) {
		return ErrorTosEmpty
	}

	return nil
}

func VerifyMailMessageDto(mailMessageDto *models.MailMessageDto) (error) {
	if nil == mailMessageDto {
		return DtoNull
	}

	err := verifyMessageDto(&mailMessageDto.MessageDto)
	if nil != err {
		return err
	}

	if common.IsStringBlank(mailMessageDto.Subject) {
		return SubjectBlank
	}

	return nil
}

func VerifySmsMessageDto(smsMessageDto *models.SmsMessageDto) (error) {
	if nil == smsMessageDto {
		return DtoNull
	}

	err := verifyMessageDto(&smsMessageDto.MessageDto)
	if nil != err {
		return err
	}

	return nil
}
