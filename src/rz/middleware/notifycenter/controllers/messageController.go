package controllers

import (
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/services"
	"rz/middleware/notifycenter/models"
)

// MVC structure
var (
	MessageController = messageController{
		SendMailControllerPack: &common.ControllerPack{
			Pattern:          "/message/send-mail",
			Method:           "POST",
			ControllerFunc:   sendMail,
			ConvertToDtoFunc: ConvertToMailMessageDto,
		},
		SendSmsControllerPack: &common.ControllerPack{
			Pattern:          "/message/send-sms",
			Method:           "POST",
			ControllerFunc:   sendSms,
			ConvertToDtoFunc: ConvertToSmsMessageDto,
		},
		QuerySmsControllerPack: &common.ControllerPack{
			Pattern:          "/message/query-sms",
			Method:           "POST",
			ControllerFunc:   querySms,
			ConvertToDtoFunc: ConvertToQuerySmsMessageRequestDtoDto,
		},
	}
)

type messageController struct {
	ControllerBase

	SendMailControllerPack *common.ControllerPack
	SendSmsControllerPack  *common.ControllerPack
	QuerySmsControllerPack *common.ControllerPack
}

func sendMail(dto interface{}) (interface{}, error) {
	mailMessageDto, ok := dto.(*models.MailMessageDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.MailMessageDto)")
	if nil != err {
		return nil, err
	}

	return services.MailMessageService.SendMail(mailMessageDto)
}

func sendSms(dto interface{}) (interface{}, error) {
	smsMessageDto, ok := dto.(*models.SmsMessageDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.SmsMessageDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsMessageService.SendSms(smsMessageDto)
}

func querySms(dto interface{}) (interface{}, error) {
	querySmsMessageRequestDto, ok := dto.(*models.QuerySmsMessageRequestDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.QuerySmsMessageRequestDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsMessageService.QuerySms(querySmsMessageRequestDto)
}
