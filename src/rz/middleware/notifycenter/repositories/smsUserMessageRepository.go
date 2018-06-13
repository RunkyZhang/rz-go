package repositories

import (
	"rz/middleware/notifycenter/models"
)

var (
	SmsUserMessageRepository smsUserMessageRepository
)

func init() {
	SmsUserMessageRepository.defaultDatabaseKey = "default"
	SmsUserMessageRepository.rawTableName = "smsUserMessage"
}

type smsUserMessageRepository struct {
	MessageRepositoryBase
}

func (myself *smsUserMessageRepository) Insert(smsTemplatePo *models.SmsUserMessagePo) (error) {
	return myself.repositoryBase.Insert(smsTemplatePo, nil)
}

func (myself *smsUserMessageRepository) SelectByPhoneNumber(nationCode string, phoneNumber string) ([]models.SmsUserMessagePo, error) {
	database, err := myself.getShardingDatabase(nil)
	if nil != err {
		return nil, err
	}

	var smsUserMessagePos []models.SmsUserMessagePo
	err = database.Where("phoneNumber=? AND nationCode=?", phoneNumber, nationCode).Find(smsUserMessagePos).Error

	return smsUserMessagePos, err
}
