package services

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/enumerations"
)

type messageServiceBase struct {
	SendChannel enumerations.SendChannel
	Prefix      string
}

func (messageServiceBase *messageServiceBase) setMessageBasePo(messageBasePo *models.MessageBasePo) {
	messageBasePo.Finished = false
	messageBasePo.SendChannel = messageServiceBase.SendChannel
	messageBasePo.States = enumerations.MessageStateToString(enumerations.Initial)
}
