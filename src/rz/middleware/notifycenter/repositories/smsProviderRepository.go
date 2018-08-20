package repositories

import (
	"rz/middleware/notifycenter/models"
	"rz/core/common"
)

var (
	SmsProviderRepository smsProviderRepository
)

func init() {
	SmsProviderRepository.DefaultDatabaseKey = "default"
	SmsProviderRepository.RawTableName = "smsProviderPo"
}

type smsProviderRepository struct {
	repositoryBase
}

func (myself *smsProviderRepository) Insert(smsProviderPo *models.SmsProviderPo) (error) {
	err := common.Assert.IsTrueToError(nil != smsProviderPo, "nil != smsProviderPo")
	if nil != err {
		return err
	}

	return myself.RepositoryBase.Insert(smsProviderPo, nil)
}

func (myself *smsProviderRepository) SelectAll() ([]*models.SmsProviderPo, error) {
	var smsProviderPos []*models.SmsProviderPo
	err := myself.RepositoryBase.SelectAll(&smsProviderPos)

	return smsProviderPos, err
}
