package services

import (
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/models"
	"time"
)

var (
	SmsUserCallbackService = smsUserCallbackService{}
)

type smsUserCallbackService struct {
}

func (*smsUserCallbackService) Add(smsUserCallbackMessageRequestExternalDto *external.SmsUserCallbackMessageRequestExternalDto) (*external.SmsUserCallbackMessageResponseExternalDto) {
	extend, err := common.StringToInt32(smsUserCallbackMessageRequestExternalDto.Extend)
	if nil != err {
		return &external.SmsUserCallbackMessageResponseExternalDto{
			Result: 1,
			Errmsg: "invalid extend",
		}
	}

	smsUserCallbackDto, err := managements.SmsUserCallbackManagement.Get(
		smsUserCallbackMessageRequestExternalDto.Nationcode,
		smsUserCallbackMessageRequestExternalDto.Mobile,
		extend)
	if nil != err {
		return &external.SmsUserCallbackMessageResponseExternalDto{
			Result: 1,
			Errmsg: "failed to get [smsUserCallbackDto]",
		}
	}

	smsUserCallbackMessageDto := models.SmsUserMessageDto{
		Content:        smsUserCallbackMessageRequestExternalDto.Text,
		Sign:           smsUserCallbackMessageRequestExternalDto.Sign,
		Time:           smsUserCallbackMessageRequestExternalDto.Time,
		CreatedTime:    time.Now().Unix(),
		Finished:       false,
		UserCallbackId: smsUserCallbackDto.Id,
	}
	smsUserCallbackMessageDto.Id = managements.SmsUserMessageManagement.BuildId(
		smsUserCallbackDto.NationCode,
		smsUserCallbackDto.PhoneNumber,
		smsUserCallbackMessageDto.CreatedTime)
	if nil == smsUserCallbackDto.UserCallbackMessages {
		smsUserCallbackDto.UserCallbackMessages = make(map[string]models.SmsUserMessageDto)
	}
	smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = smsUserCallbackMessageDto
	err = managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)
	if nil != err {
		return &external.SmsUserCallbackMessageResponseExternalDto{
			Result: 1,
			Errmsg: "failed to set [smsUserCallbackDto]",
		}
	}

	err = managements.SmsUserMessageManagement.Add(&smsUserCallbackMessageDto)
	if nil != err {
		return &external.SmsUserCallbackMessageResponseExternalDto{
			Result: 1,
			Errmsg: "failed to add [smsUserCallbackMessageDto]",
		}
	}

	return &external.SmsUserCallbackMessageResponseExternalDto{
		Result: 0,
		Errmsg: "OK",
	}
}
