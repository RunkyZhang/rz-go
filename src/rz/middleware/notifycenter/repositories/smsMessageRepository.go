package repositories

import (
	"rz/middleware/notifycenter/models"
)

type SmsMessageRepository struct {
	//reflect.ValueOf(controller)
}

func (*SmsMessageRepository) Insert(smsMessagePo *models.SmsMessagePo) error {
	return database.Create(smsMessagePo).Error
}
