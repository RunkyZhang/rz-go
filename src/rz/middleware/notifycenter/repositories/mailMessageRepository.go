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
	MailMessageRepository.DefaultDatabaseKey = "default"
	MailMessageRepository.RawTableName = "mailMessagePo"
	MailMessageRepository.GetDatabaseKeyFunc = MailMessageRepository.getDatabaseKey
	MailMessageRepository.GetTableNameFunc = MailMessageRepository.getTableName
}

func (myself *mailMessageRepository) Insert(mailMessagePo *models.MailMessagePo) (error) {
	err := common.Assert.IsNotNilToError(mailMessagePo, "mailMessagePo")
	if nil != err {
		return err
	}

	return myself.RepositoryBase.Insert(mailMessagePo, mailMessagePo.CreatedTime)
}

func (myself *mailMessageRepository) SelectById(id int, date time.Time) (*models.MailMessagePo, error) {
	mailMessagePo := &models.MailMessagePo{}

	err := myself.RepositoryBase.SelectById(id, mailMessagePo, date)

	return mailMessagePo, err
}

func (myself *mailMessageRepository) SelectByExpireTimeAndFinished(date time.Time) ([]models.MailMessagePo, error) {
	var mailMessagePos []models.MailMessagePo

	err := myself.MessageRepositoryBase.SelectByExpireTimeAndFinished(&mailMessagePos, date)

	return mailMessagePos, err
}

func (myself *mailMessageRepository) getDatabaseKey(shardParameters ...interface{}) (string) {
	return myself.DefaultDatabaseKey
}

func (myself *mailMessageRepository) getTableName(shardParameters ...interface{}) (string) {
	if nil == shardParameters || 0 == len(shardParameters) {
		return ""
	}

	date, ok := shardParameters[0].(time.Time)
	if !ok {
		return ""
	}

	return myself.RawTableName + "_" + common.Int32ToString(date.Year())
}
