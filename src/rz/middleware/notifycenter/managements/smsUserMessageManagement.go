package managements

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/common"
)

var (
	SmsUserMessageManagement = smsUserMessageManagement{}
)

func init() {
	var err error
	SmsUserMessageManagement.SendChannel = enumerations.SmsCallback
	SmsUserMessageManagement.keySuffix, err = enumerations.SendChannelToString(SmsUserMessageManagement.SendChannel)
	common.Assert.IsNilError(err, "")
	SmsUserMessageManagement.messageRepositoryBase = repositories.SmsUserMessageRepository.MessageRepositoryBase
}

type smsUserMessageManagement struct {
	MessageManagementBase
}

func (myself *smsUserMessageManagement) Add(smsUserMessagePo *models.SmsUserMessagePo) (error) {
	myself.setPoBase(&smsUserMessagePo.PoBase)
	myself.setCallbackBasePo(&smsUserMessagePo.CallbackBasePo)

	return repositories.SmsUserMessageRepository.Insert(smsUserMessagePo)
}

func (myself *smsUserMessageManagement) GetByPhoneNumber(nationCode string, phoneNumber string) ([]models.SmsUserMessagePo, error) {
	return repositories.SmsUserMessageRepository.SelectByPhoneNumber(nationCode, phoneNumber)
}

func (myself *smsUserMessageManagement) RemoveById(id string) (bool, error) {
	count, err := global.GetRedisClient().HashDelete(global.RedisKeySmsUserCallbackMessages, id)

	return 0 < count, err
}

func (myself *smsUserMessageManagement) GetAllIds() ([]string, error) {
	return global.GetRedisClient().HashKeys(global.RedisKeySmsUserCallbackMessages)
}
