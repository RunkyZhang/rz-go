package services

import (
	"time"
	"strings"

	"rz/core/common"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/global"
)

func VerifyMailMessageDto(mailMessageDto *models.MailMessageDto) (error) {
	err := common.Assert.IsTrueToError(nil != mailMessageDto, "nil != mailMessageDto")
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
	err := common.Assert.IsTrueToError(nil != smsMessageDto, "nil != smsMessageDto")
	if nil != err {
		return err
	}
	err = verifyMessageDto(&smsMessageDto.MessageBaseDto)
	if nil != err {
		return err
	}
	if nil != smsMessageDto.Parameters {
		for _, value := range smsMessageDto.Parameters {
			if strings.Contains(value, ",") {
				return exceptions.SmsParameterContainComma().AttachMessage(value)
			}
		}
	}

	return nil
}

func VerifySmsTemplateDto(smsTemplateDto *models.SmsTemplateDto) (error) {
	err := common.Assert.IsTrueToError(nil != smsTemplateDto, "nil != smsTemplateDto")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError("" != smsTemplateDto.Sign, "'' != smsTemplateDto.Sign")
	if nil != err {
		return err
	}

	if 0 >= smsTemplateDto.ContentType {
		smsTemplateDto.ContentType = 1
	}

	return nil
}

func VerifySmsProviderDto(smsProviderDto *models.SmsProviderDto) (error) {
	err := common.Assert.IsTrueToError(nil != smsProviderDto, "nil != smsProviderDto")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError("" != smsProviderDto.Id, "'' != smsProviderDto.Id")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError("" != smsProviderDto.Name, "'' != smsProviderDto.Name")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError("" != smsProviderDto.Url1, "'' != smsProviderDto.Url1")
	if nil != err {
		return err
	}

	if 0 >= smsProviderDto.Weighted {
		smsProviderDto.Weighted = 1
	}

	return nil
}

func VerifyQuerySmsUserMessagesRequestDto(querySmsUserMessagesRequestDto *models.QuerySmsUserMessagesRequestDto) (error) {
	err := common.Assert.IsTrueToError(nil != querySmsUserMessagesRequestDto, "nil != querySmsUserMessagesRequestDto")
	if nil != err {
		return err
	}
	if 0 >= querySmsUserMessagesRequestDto.SmsMessageId && "" == querySmsUserMessagesRequestDto.Content &&
		0 >= querySmsUserMessagesRequestDto.TemplateId && "" == querySmsUserMessagesRequestDto.PhoneNumber {
		return exceptions.NullQueryParameter()
	}
	if "" == querySmsUserMessagesRequestDto.NationCode {
		querySmsUserMessagesRequestDto.NationCode = global.GetConfig().SmsTencent.DefaultNationCode
	}

	return nil
}

func VerifySystemAliasPermissionDto(verifySystemAliasPermissionDto *models.SystemAliasPermissionDto) (error) {
	err := common.Assert.IsTrueToError(nil != verifySystemAliasPermissionDto, "nil != verifySystemAliasPermissionDto")
	if nil != err {
		return err
	}
	if "" == verifySystemAliasPermissionDto.SystemAlias {
		return exceptions.SystemAliasBlank()
	}

	return nil
}

func VerifyModifySystemAliasPermissionRequestDto(modifySystemAliasPermissionRequestDto *models.ModifySystemAliasPermissionRequestDto) (error) {
	err := common.Assert.IsTrueToError(nil != modifySystemAliasPermissionRequestDto, "nil != modifySystemAliasPermissionRequestDto")
	if nil != err {
		return err
	}
	if "" == modifySystemAliasPermissionRequestDto.SystemAlias {
		return exceptions.SystemAliasBlank()
	}
	if nil == modifySystemAliasPermissionRequestDto.SmsPermission &&
		nil == modifySystemAliasPermissionRequestDto.MailPermission &&
		nil == modifySystemAliasPermissionRequestDto.SmsDayFrequency &&
		nil == modifySystemAliasPermissionRequestDto.SmsHourFrequency &&
		nil == modifySystemAliasPermissionRequestDto.SmsMinuteFrequency {
		return exceptions.NullModifyParameter()
	}

	return nil
}

func verifyMessageDto(messageBaseDto *models.MessageBaseDto) (error) {
	if 0 == len(messageBaseDto.Tos) {
		return exceptions.TosEmpty()
	}
	if "" == messageBaseDto.SystemAlias {
		return exceptions.InvalidSystemAlias()
	}
	if time.Now().Unix() > messageBaseDto.ExpireTime {
		return exceptions.InvalidMessageExpireTime().AttachMessage(messageBaseDto.ExpireTime)
	}

	return nil
}
