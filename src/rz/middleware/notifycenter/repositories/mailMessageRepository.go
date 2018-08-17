package repositories

import (
	"rz/core/common"
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
	err := common.Assert.IsTrueToError(nil != mailMessagePo, "nil != mailMessagePo")
	if nil != err {
		return err
	}

	return myself.RepositoryBase.Insert(mailMessagePo, mailMessagePo.Id)
}

func (myself *mailMessageRepository) SelectById(id int64) (*models.MailMessagePo, error) {
	mailMessagePo := &models.MailMessagePo{}

	err := myself.selectById(id, mailMessagePo)

	return mailMessagePo, err
}

func (myself *mailMessageRepository) SelectByIds(ids []int64, year int) ([]*models.MailMessagePo, error) {
	var mailMessagePos []*models.MailMessagePo

	err := myself.selectByIds(ids, &mailMessagePos, year)

	return mailMessagePos, err
}

func (myself *mailMessageRepository) SelectByExpireTimeAndFinished(year int) ([]*models.MailMessagePo, error) {
	var mailMessagePos []*models.MailMessagePo

	err := myself.selectByExpireTimeAndFinished(&mailMessagePos, year)

	return mailMessagePos, err
}
