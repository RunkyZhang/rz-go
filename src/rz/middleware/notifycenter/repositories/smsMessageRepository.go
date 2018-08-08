package repositories

import (
	"time"

	"rz/middleware/notifycenter/models"
	"rz/core/common"
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

	return myself.RepositoryBase.Insert(smsMessagePo, smsMessagePo.Id)
}

func (myself *smsMessageRepository) SelectById(id int64) (*models.SmsMessagePo, error) {
	smsMessagePo := &models.SmsMessagePo{}

	err := myself.selectById(id, smsMessagePo)

	return smsMessagePo, err
}

func (myself *smsMessageRepository) SelectByIds(ids []int64, year int) ([]*models.SmsMessagePo, error) {
	var smsMessagePos []*models.SmsMessagePo

	err := myself.selectByIds(ids, &smsMessagePos, year)

	return smsMessagePos, err
}

func (myself *smsMessageRepository) SelectByExpireTimeAndFinished(year int) ([]*models.SmsMessagePo, error) {
	var smsMessagePos []*models.SmsMessagePo

	err := myself.selectByExpireTimeAndFinished(&smsMessagePos, year)

	return smsMessagePos, err
}

func (myself *smsMessageRepository) SelectByIdentifyingCode(templateId int, identifyingCode string, year int) (*models.SmsMessagePo, error) {
	database, err := myself.GetShardDatabase(year)
	if nil != err {
		return nil, err
	}

	smsMessagePo := &models.SmsMessagePo{}
	err = database.Where("templateId=? AND identifyingCode=? AND expireTime>? AND disable=0 AND deleted=0", templateId, identifyingCode, time.Now()).Last(smsMessagePo).Error

	return smsMessagePo, err
}
