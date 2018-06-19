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
	MessageRepositoryBase
}

func init() {
	SmsMessageRepository.defaultDatabaseKey = "default"
	SmsMessageRepository.rawTableName = "smsMessagePo"
	SmsMessageRepository.getDatabaseKeyFunc = SmsMessageRepository.getDatabaseKey
	SmsMessageRepository.getTableNameFunc = SmsMessageRepository.getTableName
}

func (myself *smsMessageRepository) Insert(smsMessagePo *models.SmsMessagePo) (error) {
	err := common.Assert.IsNotNilToError(smsMessagePo, "smsMessagePo")
	if nil != err {
		return err
	}

	return myself.repositoryBase.Insert(smsMessagePo, smsMessagePo.CreatedTime)
}

func (myself *smsMessageRepository) SelectById(id int, date time.Time) (*models.SmsMessagePo, error) {
	smsMessagePo := &models.SmsMessagePo{}

	err := myself.repositoryBase.SelectById(id, smsMessagePo, date)

	return smsMessagePo, err
}

func (myself *smsMessageRepository) SelectByExpireTimeAndFinished(date time.Time) ([]models.SmsMessagePo, error) {
	var smsMessagePos []models.SmsMessagePo

	err := myself.MessageRepositoryBase.SelectByExpireTimeAndFinished(&smsMessagePos, date)

	return smsMessagePos, err
}

func (myself *smsMessageRepository) SelectByIdentifyingCode(templateId int, identifyingCode string, date time.Time) (*models.SmsMessagePo, error) {
	database, err := myself.getShardingDatabase(date)
	if nil != err {
		return nil, err
	}

	smsMessagePo := &models.SmsMessagePo{}
	err = database.Where("templateId=? AND identifyingCode=? AND expireTime<? AND deleted=0", templateId, identifyingCode, time.Now()).First(smsMessagePo).Error

	return smsMessagePo, err
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
