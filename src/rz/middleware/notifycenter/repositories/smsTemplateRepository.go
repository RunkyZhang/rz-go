package repositories

import (
	"time"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
)

var (
	SmsTemplateRepository smsTemplateRepository
)

func init() {
	SmsTemplateRepository.DefaultDatabaseKey = "default"
	SmsTemplateRepository.RawTableName = "smsTemplatePo"
}

type smsTemplateRepository struct {
	common.RepositoryBase
}

func (myself *smsTemplateRepository) Insert(smsTemplatePo *models.SmsTemplatePo) (error) {
	err := common.Assert.IsNotNilToError(smsTemplatePo, "smsTemplatePo")
	if nil != err {
		return err
	}

	return myself.RepositoryBase.Insert(smsTemplatePo, nil)
}

func (myself *smsTemplateRepository) UpdateById(id int, userCallbackUrls string, pattern string) (error) {
	database, err := myself.GetShardDatabase(nil)
	if nil != err {
		return err
	}

	keyValues := map[string]interface{}{}
	keyValues["userCallbackUrls"] = userCallbackUrls
	keyValues["pattern"] = pattern
	keyValues["updatedTime"] = time.Now()

	return database.Where("id=? AND deleted=0", id).Updates(keyValues).Error
}

func (myself *smsTemplateRepository) SelectAll() ([]models.SmsTemplatePo, error) {
	var smsTemplatePos []models.SmsTemplatePo
	err := myself.RepositoryBase.SelectAll(&smsTemplatePos)

	return smsTemplatePos, err
}
