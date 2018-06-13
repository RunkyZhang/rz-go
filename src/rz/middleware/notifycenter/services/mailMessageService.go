package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
)

var (
	MailMessageService = mailMessageService{}
)

func init() {
	MailMessageService.messageManagementBase = managements.MailMessageManagement.MessageManagementBase
}

type mailMessageService struct {
	MessageServiceBase
}

func (myself *mailMessageService) SendMail(mailMessageDto *models.MailMessageDto) (int, error) {
	err := VerifyMailMessageDto(mailMessageDto)
	if nil != err {
		return 0, err
	}

	mailMessagePo := models.MailMessageDtoToPo(mailMessageDto)
	err = managements.MailMessageManagement.Add(mailMessagePo)
	if nil != err {
		return 0, err
	}

	err = managements.MailMessageManagement.EnqueueMessageIds(mailMessagePo.Id, mailMessagePo.ScheduleTime.Unix())
	if nil != err {
		myself.modifyMessageFlow(
			mailMessagePo.Id,
			&mailMessagePo.PoBase,
			&mailMessagePo.CallbackBasePo,
			enumerations.Error,
			true,
			exceptions.FailedEnqueueMessageId().AttachError(err).AttachMessage(common.Int32ToString(mailMessagePo.Id)).Error())

		return 0, err
	}

	return mailMessagePo.Id, nil
}
