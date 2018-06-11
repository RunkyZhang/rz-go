package repositories

import (
	"rz/middleware/notifycenter/models"
	"time"
)

type SmsMessageRepository struct {
	//reflect.ValueOf(controller)
}

func (*SmsMessageRepository) Insert(smsMessagePo *models.SmsMessagePo) error {
	return database.Create(smsMessagePo).Error
}

func (*SmsMessageRepository) UpdateById(id int, states string, finished bool, errorMessages string) (int64, error) {
	keyValues := map[string]interface{}{}
	keyValues["states"] = states
	keyValues["finished"] = finished
	keyValues["errorMessages"] = errorMessages
	keyValues["updatedTime"] = time.Now()

	database := database.Where("id=?", id).Updates(keyValues)
	if nil != database.Error {
		return 0, database.Error
	}

	return database.RowsAffected, nil
}

func (*SmsMessageRepository) SelectById(id int) (*models.SmsMessagePo, error) {
	smsMessagePo := &models.SmsMessagePo{}

	err := database.Where("id=? and deleted=0", id).First(smsMessagePo).Error
	if nil != err {
		return nil, err
	}

	return smsMessagePo, nil
}

func (*SmsMessageRepository) SelectByExpireTimeAndFinished() ([]models.SmsMessagePo, error) {
	var smsMessagePo []models.SmsMessagePo

	err := database.Where("finished=0 and deleted=0 and expireTime<? ", time.Now()).Find(smsMessagePo).Error
	if nil != err {
		return nil, err
	}

	return smsMessagePo, nil
}
