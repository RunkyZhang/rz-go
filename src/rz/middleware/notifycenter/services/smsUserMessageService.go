package services

import (
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/enumerations"
)

var (
	SmsUserMessageService = smsUserMessageService{}
)

func init() {
	SmsUserMessageService.messageManagementBase = managements.SmsUserMessageManagement.MessageManagementBase
}

type smsUserMessageService struct {
	MessageServiceBase
}

func (myself *smsUserMessageService) Add(
	smsUserCallbackMessageRequestExternalDto *external.SmsUserCallbackMessageRequestExternalDto) (*external.SmsUserCallbackMessageResponseExternalDto) {
	extend, err := common.StringToInt32(smsUserCallbackMessageRequestExternalDto.Extend)
	if nil != err {
		return &external.SmsUserCallbackMessageResponseExternalDto{
			Result: 1,
			Errmsg: "invalid extend",
		}
	}

	smsUserMessagePo := &models.SmsUserMessagePo{
		Content:    smsUserCallbackMessageRequestExternalDto.Text,
		Sign:       smsUserCallbackMessageRequestExternalDto.Sign,
		Time:       smsUserCallbackMessageRequestExternalDto.Time,
		NationCode: smsUserCallbackMessageRequestExternalDto.Nationcode,
		Extend:     extend,
	}

	smsTemplatePo, err := managements.SmsTemplateManagement.GetByExtend(extend)
	if nil != err {
		smsUserMessagePo.Finished = true
		smsUserMessagePo.ErrorMessages = exceptions.InvalidExtend().AttachMessage(common.Int32ToString(smsUserMessagePo.Id)).Error()
	} else {
		smsUserMessagePo.TemplateId = smsTemplatePo.Id
	}

	err = managements.SmsUserMessageManagement.Add(smsUserMessagePo)
	if nil != err {
		return &external.SmsUserCallbackMessageResponseExternalDto{
			Result: 1,
			Errmsg: exceptions.FailedAddSmsUserMessage().Error(),
		}
	}

	if false == smsUserMessagePo.Finished {
		err = managements.SmsUserMessageManagement.EnqueueMessageIds(smsUserMessagePo.Id, smsUserMessagePo.CreatedTime.Unix())
		if nil != err {
			myself.modifyMessagePo(
				&smsUserMessagePo.PoBase,
				&smsUserMessagePo.CallbackBasePo,
				enumerations.Error,
				true,
				exceptions.FailedEnqueueMessageId().AttachError(err).AttachMessage(common.Int32ToString(smsUserMessagePo.Id)).Error())

			return &external.SmsUserCallbackMessageResponseExternalDto{
				Result: 1,
				Errmsg: "Server error",
			}
		}
	}

	return &external.SmsUserCallbackMessageResponseExternalDto{
		Result: 0,
		Errmsg: "OK",
	}
}
