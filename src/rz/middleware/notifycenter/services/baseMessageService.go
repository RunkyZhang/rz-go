package services

import (
	"rz/middleware/notifycenter/models"
	"time"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/common"
)

type baseMessageService struct {
	SendChannel enumerations.SendChannel
	Prefix      string
}

func (baseMessageService *baseMessageService) setMessageDto(messageDto *models.MessageDto) (error) {
	increasing, err := managements.IncreasingManagement.Increase()
	if nil != err {
		return err
	}

	now := time.Now()
	messageDto.CreatedTime = now.Unix()
	messageDto.Id = baseMessageService.Prefix + now.Format("20060102") + common.Int64ToString(increasing)
	messageDto.SendChannel = baseMessageService.SendChannel
	messageDto.States, err = enumerations.MessageStateToString(enumerations.Initial)

	return err
}
