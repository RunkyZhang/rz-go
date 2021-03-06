package controllers

import (
	"rz/core/common"
	"rz/middleware/notifycenter/services"
	"rz/middleware/notifycenter/models"
)

// MVC structure
var (
	MessageController = messageController{
		SendMailControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/mailmessage/send",
			Method:           "POST",
			ControllerFunc:   sendMail,
			ConvertToDtoFunc: ConvertToMailMessageDto,
		},
		SendSmsControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/smsmessage/send",
			Method:           "POST",
			ControllerFunc:   sendSms,
			ConvertToDtoFunc: ConvertToSmsMessageDto,
		},
		QuerySmsControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/smsmessage/querybyids",
			Method:           "POST",
			ControllerFunc:   querySmsMessageByIds,
			ConvertToDtoFunc: ConvertToQueryMessagesByIdsRequestDto,
		},
		DisableMessageControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/smsmessage/disable",
			Method:           "POST",
			ControllerFunc:   disableSms,
			ConvertToDtoFunc: ConvertToDisableMessageRequestDto,
		},
		QuerySmsUserMessagesControllerPack: &common.ControllerPack{
			Pattern:          "/cloud.appgov.notifycenter.service/smsusermessage/query",
			Method:           "POST",
			ControllerFunc:   querySmsUserMessages,
			ConvertToDtoFunc: ConvertToQuerySmsUserMessagesRequestDto,
		},
	}
)

type messageController struct {
	ControllerBase

	SendMailControllerPack             *common.ControllerPack
	SendSmsControllerPack              *common.ControllerPack
	QuerySmsControllerPack             *common.ControllerPack
	DisableMessageControllerPack       *common.ControllerPack
	QuerySmsUserMessagesControllerPack *common.ControllerPack
}

func sendMail(dto interface{}) (interface{}, error) {
	mailMessageDto, ok := dto.(*models.MailMessageDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.MailMessageDto)")
	if nil != err {
		return nil, err
	}

	return services.MailMessageService.Send(mailMessageDto)
}

func sendSms(dto interface{}) (interface{}, error) {
	smsMessageDto, ok := dto.(*models.SmsMessageDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.SmsMessageDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsMessageService.Send(smsMessageDto)
}

func querySmsMessageByIds(dto interface{}) (interface{}, error) {
	queryMessagesByIdsRequestDto, ok := dto.(*models.QueryMessagesByIdsRequestDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.QueryMessagesByIdsRequestDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsMessageService.QueryByIds(queryMessagesByIdsRequestDto)
}

func disableSms(dto interface{}) (interface{}, error) {
	disableMessageRequestDto, ok := dto.(*models.DisableMessageRequestDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.DisableMessageRequestDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsMessageService.Disable(disableMessageRequestDto)
}

func querySmsUserMessages(dto interface{}) (interface{}, error) {
	querySmsUserMessagesRequestDto, ok := dto.(*models.QuerySmsUserMessagesRequestDto)
	err := common.Assert.IsTrueToError(ok, "dto.(*models.QuerySmsUserMessagesRequestDto)")
	if nil != err {
		return nil, err
	}

	return services.SmsUserMessageService.Query(querySmsUserMessagesRequestDto)
}
