package managements

import (
	"rz/middleware/notifycenter/global"
	"fmt"
	"rz/middleware/notifycenter/models"
	"encoding/json"
)

var (
	SmsUserMessageManagement = smsUserMessageManagement{}
)

type smsUserMessageManagement struct {
	managementBase
}

func (myself *smsUserMessageManagement) Add(smsUserMessagePo *models.SmsUserMessagePo) (error) {
	myself.setPoBase(&smsUserMessagePo.PoBase)

	return global.GetRedisClient().HashSet(global.RedisKeySmsUserCallbackMessages, smsUserMessagePo.Id, string(bytes))
}

func (myself *smsUserMessageManagement) GetById(id string) (*models.SmsUserMessageDto, error) {
	jsonString, err := global.GetRedisClient().HashGet(global.RedisKeySmsUserCallbackMessages, id)

	smsUserCallbackMessageDto := &models.SmsUserMessageDto{}
	err = json.Unmarshal([]byte(jsonString), smsUserCallbackMessageDto)
	if nil != err {
		return nil, err
	}

	return smsUserCallbackMessageDto, err
}

func (myself *smsUserMessageManagement) RemoveById(id string) (bool, error) {
	count, err := global.GetRedisClient().HashDelete(global.RedisKeySmsUserCallbackMessages, id)

	return 0 < count, err
}

func (myself *smsUserMessageManagement) GetAllIds() ([]string, error) {
	return global.GetRedisClient().HashKeys(global.RedisKeySmsUserCallbackMessages)
}

func (*smsUserMessageManagement) BuildId(nationCode string, phoneNumber string, createdTime int64) (string) {
	return fmt.Sprintf("%s_%s_%d", nationCode, phoneNumber, createdTime)
}
