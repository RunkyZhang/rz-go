package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/exceptions"
)

func VerifyMailMessageDto(mailMessageDto *models.MailMessageDto) (error) {
	err := common.Assert.IsNotNilToError(mailMessageDto, "mailMessageDto")
	if nil != err {
		return err
	}

	err = verifyMessageDto(&mailMessageDto.MessageBaseDto)
	if nil != err {
		return err
	}

	if common.IsStringBlank(mailMessageDto.Subject) {
		return exceptions.SubjectBlank()
	}

	return nil
}

func VerifySmsMessageDto(smsMessageDto *models.SmsMessageDto) (error) {
	err := common.Assert.IsNotNilToError(smsMessageDto, "smsMessageDto")
	if nil != err {
		return err
	}

	err = verifyMessageDto(&smsMessageDto.MessageBaseDto)
	if nil != err {
		return err
	}

	return nil
}

func verifyMessageDto(messageBaseDto *models.MessageBaseDto) (error) {
	if 0 == len(messageBaseDto.Tos) {
		return exceptions.ErrorTosEmpty()
	}

	return nil
}
