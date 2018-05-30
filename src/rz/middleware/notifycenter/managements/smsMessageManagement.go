package managements

import (
	"encoding/json"

	"rz/middleware/notifycenter/models"
)

var (
	SmsMessageManagement = smsMessageManagement{}
)

type smsMessageManagement struct {
	baseMessageManagement
}

func (smsMessageManagement *smsMessageManagement) AddSmsMessage(smsMessageDto *models.SmsMessageDto) (error) {
	bytes, err := json.Marshal(smsMessageDto)
	if nil != err {
		return err
	}

	return smsMessageManagement.addMessage(&smsMessageDto.BaseMessageDto, string(bytes))
}
