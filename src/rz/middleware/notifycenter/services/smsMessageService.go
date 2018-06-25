package services

import (
	"time"
	"strings"

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
	SmsMessageService.messageManagementBase = &managements.SmsMessageManagement.MessageManagementBase
}

type smsMessageService struct {
	MessageServiceBase
}

func (myself *smsMessageService) SendSms(smsMessageDto *models.SmsMessageDto) (int, error) {
	err := common.Assert.IsNotNilToError(smsMessageDto, "smsMessageDto")
	if nil != err {
		return 0, err
	}

	err = VerifySmsMessageDto(smsMessageDto)
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

func (myself *smsMessageService) QuerySms(querySmsMessageRequestDto *models.QuerySmsMessageRequestDto) (*models.SmsMessageDto, error) {
	smsMessagePo, err := managements.SmsMessageManagement.GetById(querySmsMessageRequestDto.Id, time.Unix(querySmsMessageRequestDto.CreatedTime, 0))
	if nil != err {
		return nil, err
	}

	if !strings.EqualFold(smsMessagePo.SystemAlias, querySmsMessageRequestDto.SystemAlias) {
		return nil, exceptions.MessageSystemAliasMotMatch().AttachMessage(querySmsMessageRequestDto.SystemAlias)
	}

	return models.SmsMessagePoToDto(smsMessagePo), err
}
