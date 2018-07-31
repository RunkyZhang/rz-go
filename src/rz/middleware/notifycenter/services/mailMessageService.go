package services

import (
	"time"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
)

var (
	MailMessageService = mailMessageService{}
)

func init() {
	MailMessageService.messageManagementBase = &managements.MailMessageManagement.MessageManagementBase
}

type mailMessageService struct {
	MessageServiceBase
}

func (myself *mailMessageService) Send(mailMessageDto *models.MailMessageDto) (int64, error) {
	err := VerifyMailMessageDto(mailMessageDto)
	if nil != err {
		return 0, err
	}
	systemAliasPermissionPo, err := managements.SystemAliasPermissionManagement.GetById(mailMessageDto.SystemAlias)
	if nil != err || 0 == systemAliasPermissionPo.SmsPermission {
		return 0, exceptions.NotSendMailPermission().AttachError(err).AttachMessage(mailMessageDto.SystemAlias)
	}

	mailMessagePo := models.MailMessageDtoToPo(mailMessageDto)
	mailMessagePo.CreatedTime = time.Now()
	mailMessagePo.Id, err = managements.MailMessageManagement.GenerateId(mailMessagePo.CreatedTime.Year())
	if nil != err {
		return 0, exceptions.FailedGenerateMessageId().AttachError(err)
	}

	err = managements.MailMessageManagement.Add(mailMessagePo)
	if nil != err {
		return 0, err
	}

	err = managements.MailMessageManagement.EnqueueIds(mailMessagePo.Id, mailMessagePo.ScheduleTime.Unix())
	if nil != err {
		now := time.Now()
		finished := true
		managements.ModifyMessageFlowAsync(
			myself.messageManagementBase,
			mailMessagePo.Id,
			enumerations.Initial,
			enumerations.Error,
			exceptions.FailedEnqueueMessageId().AttachError(err).AttachMessage(mailMessagePo.Id).Error(),
			&finished,
			&now,
			mailMessagePo.CreatedTime.Year())

		return 0, err
	}

	return mailMessagePo.Id, nil
}
