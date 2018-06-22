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
	err := common.Assert.IsNotNilToError(mailMessagePo, "mailMessagePo")
	if nil != err {
		return err
	}

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

func (myself *mailMessageRepository) getDatabaseKey(shardParameters ...interface{}) (string) {
	return myself.defaultDatabaseKey
}

func (myself *mailMessageRepository) getTableName(shardParameters ...interface{}) (string) {
	if nil == shardParameters || 0 == len(shardParameters) {
		return ""
	}

	date, ok := shardParameters[0].(time.Time)
	if !ok {
		return ""
	}

	return myself.rawTableName + "_" + common.Int32ToString(date.Year())
}
