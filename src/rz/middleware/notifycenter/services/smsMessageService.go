package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
	"time"
)

var (
	SmsMessageService = smsMessageService{}
)

func init() {
	SmsMessageService.messageManagementBase = &managements.SmsMessageManagement.MessageManagementBase
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
		managements.ModifyMessageFlowAsync(
			myself.messageManagementBase,
			smsMessagePo.Id,
			&smsMessagePo.PoBase,
			&smsMessagePo.CallbackBasePo,
			enumerations.Error,
			true,
			time.Now(),
			exceptions.FailedEnqueueMessageId().AttachError(err).AttachMessage(common.Int32ToString(smsMessagePo.Id)).Error())

		return 0, err
	}

	return smsMessagePo.Id, err
}
