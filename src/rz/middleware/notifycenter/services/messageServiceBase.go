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

func (myself *MessageServiceBase) modifyMessageFlow(
	messageId int,
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
		messageId,
		callbackBasePo.States,
		finished,
		errorMessages,
		poBase.CreatedTime)
	if nil != err || 0 == affectedCount {
		fmt.Printf("failed to modify message(%d) state. error: %s", messageId, err)
	}
}
