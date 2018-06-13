package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
)

var (
	SmsMessageService = smsMessageService{}
)

func init() {
	SmsMessageService.messageManagementBase = managements.SmsMessageManagement.MessageManagementBase
}

type smsMessageService struct {
	MessageServiceBase
}

func (myself *smsMessageService) SendSms(smsMessageDto *models.SmsMessageDto) (int, error) {
	err := VerifySmsMessageDto(smsMessageDto)
	if nil != err {
		return 0, err
	}

	smsMessagePo := models.SmsMessageDtoToPo(smsMessageDto)
	err = managements.SmsMessageManagement.Add(smsMessagePo)
	if nil != err {
		return 0, err
	}

	err = managements.SmsMessageManagement.EnqueueMessageIds(smsMessagePo.Id, smsMessagePo.ScheduleTime.Unix())
	if nil != err {
		myself.modifyMessagePo(
			&smsMessagePo.PoBase,
			&smsMessagePo.CallbackBasePo,
			enumerations.Error,
			true,
			exceptions.FailedEnqueueMessageId().AttachError(err).AttachMessage(common.Int32ToString(smsMessagePo.Id)).Error())

		return 0, err
	}

	return smsMessagePo.Id, err
}
