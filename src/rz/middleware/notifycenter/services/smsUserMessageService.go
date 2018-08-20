package services

import (
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/managements"
	"rz/core/common"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/enumerations"
	"time"
	"rz/middleware/notifycenter/provider"
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

func (myself *smsUserMessageService) TencentCallback(tencentSmsUserCallbackRequestDto *external.TencentSmsUserCallbackRequestDto) (*external.TencentSmsUserCallbackResponseDto) {
	err := common.Assert.IsTrueToError(nil != tencentSmsUserCallbackRequestDto, "nil != tencentSmsUserCallbackRequestDto")
	if nil != err {
		return &external.TencentSmsUserCallbackResponseDto{
			Result: 1,
			Errmsg: err.Error(),
		}
	}
	extend, err := common.StringToInt32(tencentSmsUserCallbackRequestDto.Extend)
	if nil != err {
		return &external.TencentSmsUserCallbackResponseDto{
			Result: 1,
			Errmsg: err.Error(),
		}
	}

	err = myself.callback(
		tencentSmsUserCallbackRequestDto.Mobile,
		tencentSmsUserCallbackRequestDto.Text,
		extend,
		tencentSmsUserCallbackRequestDto.Time,
		tencentSmsUserCallbackRequestDto.Nationcode,
		tencentSmsUserCallbackRequestDto.Sign,
		provider.SmsTencentProvider.Id)
	if nil != err {
		return &external.TencentSmsUserCallbackResponseDto{
			Result: 1,
			Errmsg: err.Error(),
		}
	}

	return &external.TencentSmsUserCallbackResponseDto{
		Result: 0,
		Errmsg: "OK",
	}
}

func (myself *smsUserMessageService) DahanCallbacks(dahanSmsUserCallbackRequestDto *external.DahanSmsUserCallbackRequestDto) (*external.DahanSmsUserCallbackResponseDto) {
	err := common.Assert.IsTrueToError(nil != dahanSmsUserCallbackRequestDto, "nil != dahanSmsUserCallbackRequestDto")
	if nil != err {
		return &external.DahanSmsUserCallbackResponseDto{
			Status: err.Error(),
		}
	}

	if nil != dahanSmsUserCallbackRequestDto.Delivers {
		for _, deliver := range dahanSmsUserCallbackRequestDto.Delivers {
			err = myself.dahanCallback(deliver)
			if nil != err {
				common.GetLogging().Warn(err, "Failed to save callback message")
			}
		}
	}

	return &external.DahanSmsUserCallbackResponseDto{
		Status: "success",
	}
}

func (myself *smsUserMessageService) dahanCallback(dahanSmsUserCallbackDeliverRequestDto *external.DahanSmsUserCallbackDeliverRequestDto) (error) {
	err := common.Assert.IsTrueToError(nil != dahanSmsUserCallbackDeliverRequestDto, "nil != dahanSmsUserCallbackDeliverRequestDto")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError(5 <= len(dahanSmsUserCallbackDeliverRequestDto.SubCode), "5 <= len(dahanSmsUserCallbackDeliverRequestDto.SubCode")
	if nil != err {
		return err
	}
	extend := -1
	sign := dahanSmsUserCallbackDeliverRequestDto.SubCode[0:4]
	if 5 < len(dahanSmsUserCallbackDeliverRequestDto.SubCode) {
		extend, err = common.StringToInt32(dahanSmsUserCallbackDeliverRequestDto.SubCode[5:])
		if nil != err {
			return err
		}
	}
	dateTime, err := time.Parse("2006-01-02 15:04:05", dahanSmsUserCallbackDeliverRequestDto.DeliverTime)
	if nil != err {
		return err
	}

	// {"result":"0","desc":"成功","delivers":[{"phone":"13818530040","content":"刚刚","subcode":"566013333","delivertime":"2018-08-15 14:10:08"}]}
	err = myself.callback(
		dahanSmsUserCallbackDeliverRequestDto.Phone,
		dahanSmsUserCallbackDeliverRequestDto.Content,
		extend,
		dateTime.Unix(),
		"86",
		sign,
		provider.SmsDahanProvider.Id)
	if nil != err {
		return err
	}

	return nil
}

func (myself *smsUserMessageService) callback(phoneNumber string, context string, extend int, dateTime int64, nationCode string, sign string, fromProviderId string) (error) {
	smsUserMessagePo := &models.SmsUserMessagePo{
		Content:        context,
		Sign:           sign,
		Time:           dateTime,
		NationCode:     nationCode,
		PhoneNumber:    phoneNumber,
		Extend:         extend,
		FromProviderId: fromProviderId,
	}

	smsUserMessagePo.ExpireTime = time.Now().Add(7 * 24 * time.Hour)
	smsTemplatePo, err := managements.SmsTemplateManagement.GetByExtend(extend)
	if nil != err {
		smsUserMessagePo.Finished = true
		smsUserMessagePo.ErrorMessages = exceptions.InvalidExtend().AttachMessage(smsUserMessagePo.Id).Error()
	} else {
		smsUserMessagePo.TemplateId = smsTemplatePo.Id
	}
	smsUserMessagePo.CreatedTime = time.Now()
	smsUserMessagePo.Id, err = managements.SmsUserMessageManagement.GenerateId(smsUserMessagePo.CreatedTime.Year())
	if nil != err {
		return exceptions.FailedGenerateMessageId().AttachError(err)
	}

	err = managements.SmsUserMessageManagement.Add(smsUserMessagePo)
	if nil != err {
		return exceptions.FailedAddSmsUserMessage().AttachError(err)
	}

	if false == smsUserMessagePo.Finished {
		err = managements.SmsUserMessageManagement.EnqueueIds(smsUserMessagePo.Id, smsUserMessagePo.CreatedTime.Unix())
		if nil != err {
			now := time.Now()
			finished := true
			managements.ModifyMessageFlowAsync(
				myself.messageManagementBase,
				smsUserMessagePo.Id,
				enumerations.Initial,
				enumerations.Error,
				exceptions.FailedEnqueueMessageId().AttachError(err).AttachMessage(smsUserMessagePo.Id).Error(),
				"",
				&finished,
				&now,
				smsUserMessagePo.CreatedTime.Year())

			return exceptions.InternalServerError().AttachError(err)
		}
	}

	return nil
}

func (myself *smsUserMessageService) Query(querySmsUserMessagesRequestDto *models.QuerySmsUserMessagesRequestDto) ([]*models.SmsUserMessageDto, error) {
	err := VerifyQuerySmsUserMessagesRequestDto(querySmsUserMessagesRequestDto)
	if nil != err {
		return nil, err
	}

	smsUserMessagePos, err := managements.SmsUserMessageManagement.Query(
		querySmsUserMessagesRequestDto.SmsMessageId,
		querySmsUserMessagesRequestDto.Content,
		querySmsUserMessagesRequestDto.NationCode,
		querySmsUserMessagesRequestDto.PhoneNumber,
		querySmsUserMessagesRequestDto.TemplateId,
		querySmsUserMessagesRequestDto.Year)

	return models.SmsUserMessagePosToDtos(smsUserMessagePos), err
}
