package repositories

import (
	"time"
	"rz/middleware/notifycenter/models"
	"rz/core/common"
)

var (
	SmsTemplateRepository smsTemplateRepository
)

func init() {
	SmsTemplateRepository.DefaultDatabaseKey = "default"
	SmsTemplateRepository.RawTableName = "smsTemplatePo"
}

type smsTemplateRepository struct {
	repositoryBase
}

func (myself *smsTemplateRepository) Insert(smsTemplatePo *models.SmsTemplatePo) (error) {
	err := common.Assert.IsTrueToError(nil != smsTemplatePo, "nil != smsTemplatePo")
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

func (myself *smsTemplateRepository) SelectAll() ([]*models.SmsTemplatePo, error) {
	var smsTemplatePos []*models.SmsTemplatePo
	err := myself.RepositoryBase.SelectAll(&smsTemplatePos)

	return smsTemplatePos, err
}

func (myself *smsTemplateRepository) CountByExtend(extend int) (int64, error) {
	database, err := myself.GetShardDatabase(nil)
	if nil != err {
		return 0, err
	}

	var count int64
	err = database.Where("extend=? AND deleted=0", extend).Count(&count).Error

	return count, err
}
