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

func (baseMessageService *baseMessageService) setMessageDto(baseMessageDto *models.BaseMessageDto) (error) {
	increasing, err := managements.IncreasingManagement.Increase()
	if nil != err {
		return err
	}

	now := time.Now()
	baseMessageDto.Finished = false
	baseMessageDto.CreatedTime = now.Unix()
	baseMessageDto.Id = baseMessageService.Prefix + now.Format("20060102") + common.Int64ToString(increasing)
	baseMessageDto.SendChannel = baseMessageService.SendChannel
	baseMessageDto.States = enumerations.MessageStateToString(enumerations.Initial)

	return err
}
