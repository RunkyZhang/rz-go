package repositories

import (
	"rz/middleware/notifycenter/models"
	"time"
	"rz/middleware/notifycenter/common"
)

var (
	SmsMessageRepository smsMessageRepository
)

type smsMessageRepository struct {
	messageRepositoryBase
}

func init() {
	SmsMessageRepository.defaultDatabaseKey = "default"
	SmsMessageRepository.rawTableName = "smsMessagePo"
	SmsMessageRepository.getDatabaseKeyFunc = SmsMessageRepository.getDatabaseKey
	SmsMessageRepository.getTableNameFunc = SmsMessageRepository.getTableName
}

func (myself *smsMessageRepository) Insert(smsMessagePo *models.SmsMessagePo) (error) {
	return myself.repositoryBase.Insert(smsMessagePo, smsMessagePo.CreatedTime)
}

func (myself *smsMessageRepository) SelectById(id int, date time.Time) (*models.SmsMessagePo, error) {
	smsMessagePo := &models.SmsMessagePo{}

	err := myself.repositoryBase.SelectById(id, smsMessagePo, date)
	if nil != err {
		return nil, err
	}

	return smsMessagePo, nil
}

func (myself *smsMessageRepository) SelectByExpireTimeAndFinished(date time.Time) ([]models.SmsMessagePo, error) {
	var smsMessagePos []models.SmsMessagePo

	err := myself.messageRepositoryBase.SelectByExpireTimeAndFinished(smsMessagePos, date)
	if nil != err {
		return nil, err
	}

	return smsMessagePos, nil
}

func (myself *smsMessageRepository) getDatabaseKey(shardingParameters ...interface{}) (string) {
	return myself.defaultDatabaseKey
}

func (myself *smsMessageRepository) getTableName(shardingParameters ...interface{}) (string) {
	if nil == shardingParameters || 0 == len(shardingParameters) {
		return ""
	}

	date, ok := shardingParameters[0].(time.Time)
	if !ok {
		return ""
	}

	return myself.rawTableName + "_" + common.Int32ToString(date.Year())
}
