package services

import (
	"time"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/enumerations"
)

type messageServiceBase struct {
	SendChannel enumerations.SendChannel
	Prefix      string
}

func (messageServiceBase *messageServiceBase) setMessageBasePo(messageBasePo *models.MessageBasePo) {
	now := time.Now()
	messageBasePo.Finished = false
	messageBasePo.CreatedTime = now
	messageBasePo.UpdatedTime = now
	messageBasePo.Deleted = false
	messageBasePo.SendChannel = messageServiceBase.SendChannel
	messageBasePo.States = enumerations.MessageStateToString(enumerations.Initial)
}
