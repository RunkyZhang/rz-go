package repositories

import (
	"time"

	"rz/middleware/notifycenter/models"
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

	return myself.RepositoryBase.Insert(smsUserMessagePo, smsUserMessagePo.Id)
}

func (myself *smsUserMessageRepository) UpdateSmsMessageIdById(id int64, smsMessageId int64) (int64, error) {
	database, err := myself.GetShardDatabase(id)
	if nil != err {
		return 0, err
	}

	keyValues := map[string]interface{}{}
	keyValues["smsMessageId"] = smsMessageId
	keyValues["updatedTime"] = time.Now()
	database = database.Where("id=?", id).Updates(keyValues)

	return database.RowsAffected, database.Error
}

func (myself *smsUserMessageRepository) SelectById(id int64) (*models.SmsUserMessagePo, error) {
	smsUserMessagePo := &models.SmsUserMessagePo{}

	err := myself.selectById(id, smsUserMessagePo)

	return smsUserMessagePo, err
}

func (myself *smsUserMessageRepository) Select(smsMessageId int64, content string, nationCode string, phoneNumber string, templateId int, year int) ([]*models.SmsUserMessagePo, error) {
	database, err := myself.GetShardDatabase(year)
	if nil != err {
		return nil, err
	}

	var parameters []interface{}
	whereSql := ""
	if 0 < smsMessageId {
		whereSql += "`smsMessageId`=?"
		parameters = append(parameters, smsMessageId)
	}
	if "" != content {
		if "" != whereSql {
			whereSql += " AND "
		}

		whereSql += "`content`=?"
		parameters = append(parameters, content)
	}
	if "" != nationCode {
		if "" != whereSql {
			whereSql += " AND "
		}

		whereSql += "`nationCode`=?"
		parameters = append(parameters, nationCode)
	}
	if "" != phoneNumber {
		if "" != whereSql {
			whereSql += " AND "
		}

		whereSql += "`phoneNumber`=?"
		parameters = append(parameters, phoneNumber)
	}
	if 0 < templateId {
		if "" != whereSql {
			whereSql += " AND "
		}

		whereSql += "`templateId`=?"
		parameters = append(parameters, templateId)
	}

	var smsUserMessagePos []*models.SmsUserMessagePo
	err = database.Where(whereSql, parameters...).Find(&smsUserMessagePos).Error

	return smsUserMessagePos, err
}

func (myself *smsUserMessageRepository) SelectByPhoneNumber(nationCode string, phoneNumber string, year int) ([]*models.SmsUserMessagePo, error) {
	database, err := myself.GetShardDatabase(year)
	if nil != err {
		return nil, err
	}

	var smsUserMessagePos []*models.SmsUserMessagePo
	err = database.Where("phoneNumber=? AND nationCode=?", phoneNumber, nationCode).Find(&smsUserMessagePos).Error

	return smsUserMessagePos, err
}
