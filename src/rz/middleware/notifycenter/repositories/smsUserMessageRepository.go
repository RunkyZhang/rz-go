package repositories

import (
	"rz/middleware/notifycenter/models"
	"time"
	"rz/middleware/notifycenter/common"
)

var (
	SmsUserMessageRepository smsUserMessageRepository
)

func init() {
	SmsUserMessageRepository.DefaultDatabaseKey = "default"
	SmsUserMessageRepository.RawTableName = "smsUserMessagePo"
	SmsUserMessageRepository.GetDatabaseKeyFunc = SmsUserMessageRepository.getDatabaseKey
	SmsUserMessageRepository.GetTableNameFunc = SmsUserMessageRepository.getTableName
}

type smsUserMessageRepository struct {
	MessageRepositoryBase
}

func (myself *smsUserMessageRepository) Insert(smsUserMessagePo *models.SmsUserMessagePo) (error) {
	err := common.Assert.IsNotNilToError(smsUserMessagePo, "smsUserMessagePo")
	if nil != err {
		return err
	}

	return myself.RepositoryBase.Insert(smsUserMessagePo, smsUserMessagePo.CreatedTime)
}

func (myself *smsUserMessageRepository) SelectById(id int, date time.Time) (*models.SmsUserMessagePo, error) {
	smsUserMessagePo := &models.SmsUserMessagePo{}

	err := myself.RepositoryBase.SelectById(id, smsUserMessagePo, date)

	return smsUserMessagePo, err
}

func (myself *smsUserMessageRepository) SelectByPhoneNumber(nationCode string, phoneNumber string) ([]models.SmsUserMessagePo, error) {
	database, err := myself.GetShardDatabase(nil)
	if nil != err {
		return nil, err
	}

	var smsUserMessagePos []models.SmsUserMessagePo
	err = database.Where("phoneNumber=? AND nationCode=?", phoneNumber, nationCode).Find(smsUserMessagePos).Error

	return smsUserMessagePos, err
}

func (myself *smsUserMessageRepository) getDatabaseKey(shardParameters ...interface{}) (string) {
	return myself.DefaultDatabaseKey
}

func (myself *smsUserMessageRepository) getTableName(shardParameters ...interface{}) (string) {
	if nil == shardParameters || 0 == len(shardParameters) {
		return ""
	}

	date, ok := shardParameters[0].(time.Time)
	if !ok {
		return ""
	}

	return myself.RawTableName + "_" + common.Int32ToString(date.Year())
}
