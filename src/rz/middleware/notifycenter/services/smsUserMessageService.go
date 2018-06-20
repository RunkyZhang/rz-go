package services

import (
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/enumerations"
	"time"
)

var (
	SmsUserMessageService = smsUserMessageService{}
)

func init() {
	SmsUserMessageService.messageManagementBase = &managements.SmsUserMessageManagement.MessageManagementBase
}

type smsUserMessageService struct {
	MessageServiceBase
}

func (myself *smsUserMessageService) Callback(smsUserCallbackMessageRequestExternalDto *external.SmsUserCallbackMessageRequestExternalDto) (*external.SmsUserCallbackMessageResponseExternalDto) {
	err := common.Assert.IsNotNilToError(smsUserCallbackMessageRequestExternalDto, "smsUserCallbackMessageRequestExternalDto")
	if nil != err {
		return &external.SmsUserCallbackMessageResponseExternalDto{
			Result: 1,
			Errmsg: "invalid request body",
		}
	}
	extend, err := common.StringToInt32(smsUserCallbackMessageRequestExternalDto.Extend)
	if nil != err {
		return &external.SmsUserCallbackMessageResponseExternalDto{
			Result: 1,
			Errmsg: "invalid extend",
		}
	}

	smsUserMessagePo := &models.SmsUserMessagePo{
		Content:     smsUserCallbackMessageRequestExternalDto.Text,
		Sign:        smsUserCallbackMessageRequestExternalDto.Sign,
		Time:        smsUserCallbackMessageRequestExternalDto.Time,
		NationCode:  smsUserCallbackMessageRequestExternalDto.Nationcode,
		PhoneNumber: smsUserCallbackMessageRequestExternalDto.Mobile,
		Extend:      extend,
	}
	smsUserMessagePo.ExpireTime = time.Now().Add(7 * 24 * time.Hour)

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
			Errmsg: exceptions.FailedAddSmsUserMessage().AttachError(err).Error(),
		}
	}

	if false == smsUserMessagePo.Finished {
		err = managements.SmsUserMessageManagement.EnqueueMessageIds(smsUserMessagePo.Id, smsUserMessagePo.CreatedTime.Unix())
		if nil != err {
			managements.ModifyMessageFlowAsync(
				myself.messageManagementBase,
				smsUserMessagePo.Id,
				&smsUserMessagePo.PoBase,
				&smsUserMessagePo.CallbackBasePo,
				enumerations.Error,
				true,
				time.Now(),
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
