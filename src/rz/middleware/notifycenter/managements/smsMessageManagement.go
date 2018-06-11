package managements

import (

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
)

var (
	SmsMessageManagement = smsMessageManagement{}
)

type smsMessageManagement struct {
	messageManagementBase

	smsMessageRepository *repositories.SmsMessageRepository
}

func (smsMessageManagement *smsMessageManagement) Add(smsMessagePo *models.SmsMessagePo) (error) {
	//bytes, err := json.Marshal(smsMessageDto)
	//if nil != err {
	//	return err
	//}
	//
	//return smsMessageManagement.addMessage(&smsMessageDto.MessageBaseDto, string(bytes))

	return smsMessageManagement.smsMessageRepository.Insert(smsMessagePo)
}
