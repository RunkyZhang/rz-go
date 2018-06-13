package repositories

import (
	"time"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/models"
)

var (
	MailMessageRepository mailMessageRepository
)

type mailMessageRepository struct {
	MessageRepositoryBase
}

func init() {
	MailMessageRepository.defaultDatabaseKey = "default"
	MailMessageRepository.rawTableName = "mailMessagePo"
	MailMessageRepository.getDatabaseKeyFunc = MailMessageRepository.getDatabaseKey
	MailMessageRepository.getTableNameFunc = MailMessageRepository.getTableName
}

func (myself *mailMessageRepository) Insert(mailMessagePo *models.MailMessagePo) (error) {
	return myself.repositoryBase.Insert(mailMessagePo, mailMessagePo.CreatedTime)
}

func (myself *mailMessageRepository) SelectById(id int, date time.Time) (*models.MailMessagePo, error) {
	mailMessagePo := &models.MailMessagePo{}

	err := myself.repositoryBase.SelectById(id, mailMessagePo, date)

	return mailMessagePo, err
}

func (myself *mailMessageRepository) SelectByExpireTimeAndFinished(date time.Time) ([]models.MailMessagePo, error) {
	var mailMessagePos []models.MailMessagePo

	err := myself.MessageRepositoryBase.SelectByExpireTimeAndFinished(&mailMessagePos, date)

	return mailMessagePos, err
}

func (myself *mailMessageRepository) getDatabaseKey(shardingParameters ...interface{}) (string) {
	return myself.defaultDatabaseKey
}

func (myself *mailMessageRepository) getTableName(shardingParameters ...interface{}) (string) {
	if nil == shardingParameters || 0 == len(shardingParameters) {
		return ""
	}

	date, ok := shardingParameters[0].(time.Time)
	if !ok {
		return ""
	}

	return myself.rawTableName + "_" + common.Int32ToString(date.Year())
}
