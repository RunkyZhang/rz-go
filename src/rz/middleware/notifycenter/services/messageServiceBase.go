package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"fmt"
	"rz/middleware/notifycenter/enumerations"
)

type MessageServiceBase struct {
	messageManagementBase managements.MessageManagementBase
}

func (myself *MessageServiceBase) modifyMessagePo(
	poBase *models.PoBase,
	callbackBasePo *models.CallbackBasePo,
	messageState enumerations.MessageState,
	finished bool,
	errorMessage string) {
	state := enumerations.MessageStateToString(messageState)
	callbackBasePo.States = callbackBasePo.States + "+" + state
	var errorMessages string
	if "" == errorMessage {
		errorMessages = ""
	} else {
		callbackBasePo.ErrorMessages = callbackBasePo.ErrorMessages + "+++" + errorMessage
		errorMessages = callbackBasePo.ErrorMessages
	}

	affectedCount, err := myself.messageManagementBase.ModifyById(
		poBase.Id,
		callbackBasePo.States,
		finished,
		errorMessages,
		poBase.CreatedTime)
	if nil != err || 0 == affectedCount {
		fmt.Println("failed to modify message(", poBase.Id, ") state. error: ", err.Error())
	}
}
