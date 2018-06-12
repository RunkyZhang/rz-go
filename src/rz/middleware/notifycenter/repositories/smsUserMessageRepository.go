package repositories

import (
	"time"
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
	repositoryBase
}

func (myself *smsUserMessageRepository) Insert(smsTemplatePo *models.SmsUserMessagePo) (error) {
	return myself.repositoryBase.Insert(smsTemplatePo, nil)
}

func (myself *smsUserMessageRepository) UpdateById(id int, userCallbackUrls string, pattern string) (error) {
	database, err := myself.getDatabase(nil)
	if nil != err {
		return err
	}
	tableName := myself.getTableName(nil)

	keyValues := map[string]interface{}{}
	keyValues["userCallbackUrls"] = userCallbackUrls
	keyValues["pattern"] = pattern
	keyValues["updatedTime"] = time.Now()

	return database.Table(tableName).Where("id=?", id).Updates(keyValues).Error
}
