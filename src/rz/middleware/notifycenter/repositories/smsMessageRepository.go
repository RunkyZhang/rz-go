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
	SmsMessageRepository.DefaultDatabaseKey = "default"
	SmsMessageRepository.RawTableName = "smsMessagePo"
	SmsMessageRepository.GetDatabaseKeyFunc = SmsMessageRepository.getDatabaseKey
	SmsMessageRepository.GetTableNameFunc = SmsMessageRepository.getTableName
}

func (myself *smsMessageRepository) Insert(smsMessagePo *models.SmsMessagePo) (error) {
	err := common.Assert.IsNotNilToError(smsMessagePo, "smsMessagePo")
	if nil != err {
		return err
	}

	return myself.RepositoryBase.Insert(smsMessagePo, smsMessagePo.CreatedTime)
}

func (myself *smsMessageRepository) SelectById(id int, date time.Time) (*models.SmsMessagePo, error) {
	smsMessagePo := &models.SmsMessagePo{}

	err := myself.RepositoryBase.SelectById(id, smsMessagePo, date)

	return smsMessagePo, err
}

func (myself *smsMessageRepository) SelectByExpireTimeAndFinished(date time.Time) ([]models.SmsMessagePo, error) {
	var smsMessagePos []models.SmsMessagePo

	err := myself.MessageRepositoryBase.SelectByExpireTimeAndFinished(&smsMessagePos, date)

	return smsMessagePos, err
}

func (myself *smsMessageRepository) SelectByIdentifyingCode(templateId int, identifyingCode string, date time.Time) (*models.SmsMessagePo, error) {
	database, err := myself.GetShardDatabase(date)
	if nil != err {
		return nil, err
	}

	smsMessagePo := &models.SmsMessagePo{}
	err = database.Where("templateId=? AND identifyingCode=? AND expireTime>? AND deleted=0", templateId, identifyingCode, time.Now()).First(smsMessagePo).Error

	return smsMessagePo, err
}

func (myself *smsMessageRepository) getDatabaseKey(shardParameters ...interface{}) (string) {
	return myself.DefaultDatabaseKey
}

func (myself *smsMessageRepository) getTableName(shardParameters ...interface{}) (string) {
	if nil == shardParameters || 0 == len(shardParameters) {
		return ""
	}

	date, ok := shardParameters[0].(time.Time)
	if !ok {
		return ""
	}

	return myself.RawTableName + "_" + common.Int32ToString(date.Year())
}
