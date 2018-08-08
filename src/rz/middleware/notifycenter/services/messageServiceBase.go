package services

import (
	"rz/middleware/notifycenter/managements"
	"strings"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/core/common"
	"rz/middleware/notifycenter/enumerations"
)

type MessageServiceBase struct {
	messageManagementBase *managements.MessageManagementBase
}

func (myself *MessageServiceBase) Disable(disableMessageRequestDto *models.DisableMessageRequestDto) (bool, error) {
	err := common.Assert.IsNotNilToError(disableMessageRequestDto, "disableMessageRequestDto")
	if nil != err {
		return false, err
	}
	smsMessagePo, err := managements.SmsMessageManagement.GetById(disableMessageRequestDto.Id)
	if nil != err {
		return false, err
	}
	if !strings.EqualFold(smsMessagePo.SystemAlias, disableMessageRequestDto.SystemAlias) {
		return false, exceptions.MessageSystemAliasMotMatch().AttachMessage(disableMessageRequestDto.SystemAlias)
	}
	if enumerations.MessageStateToString(enumerations.Initial) != smsMessagePo.States {
		return false, exceptions.MessageNotInitialState().AttachMessage(smsMessagePo.States)
	}

	rowsAffected, err := managements.SmsMessageManagement.Disable(disableMessageRequestDto.Id)
	return 0 < rowsAffected, err
}
