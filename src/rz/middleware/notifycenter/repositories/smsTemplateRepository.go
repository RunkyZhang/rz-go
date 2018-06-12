package repositories

import (
	"time"
	"rz/middleware/notifycenter/models"
)

var (
	SmsTemplateRepository smsTemplateRepository
)

func init() {
	SmsTemplateRepository.defaultDatabaseKey = "default"
	SmsTemplateRepository.rawTableName = "smsTemplatePo"
}

type smsTemplateRepository struct {
	repositoryBase
}

func (myself *smsTemplateRepository) Insert(smsTemplatePo *models.SmsTemplatePo) (error) {
	return myself.repositoryBase.Insert(smsTemplatePo, nil)
}

func (myself *smsTemplateRepository) UpdateById(id int, userCallbackUrls string, pattern string) (error) {
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

func (myself *smsTemplateRepository) SelectAll() ([]models.SmsTemplatePo, error) {
	var smsTemplatePos []models.SmsTemplatePo
	err := myself.repositoryBase.SelectAll(smsTemplatePos)
	if nil != err {
		return nil, err
	}

	return smsTemplatePos, nil
}
