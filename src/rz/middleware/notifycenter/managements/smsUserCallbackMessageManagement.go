package managements

import (
	"rz/middleware/notifycenter/global"
	"fmt"
	"rz/middleware/notifycenter/models"
	"encoding/json"
)

var (
	SmsUserCallbackMessageManagement = smsUserCallbackMessageManagement{}
)

type smsUserCallbackMessageManagement struct {
}

func (smsUserCallbackMessageManagement *smsUserCallbackMessageManagement) Add(smsUserCallbackMessageDto *models.SmsUserCallbackMessageDto) (error) {
	bytes, err := json.Marshal(smsUserCallbackMessageDto)
	if nil != err {
		return err
	}

	return global.GetRedisClient().HashSet(global.RedisKeySmsUserCallbackMessages, smsUserCallbackMessageDto.Id, string(bytes))
}

func (smsUserCallbackMessageManagement *smsUserCallbackMessageManagement) GetById(id string) (*models.SmsUserCallbackMessageDto, error) {
	jsonString, err := global.GetRedisClient().HashGet(global.RedisKeySmsUserCallbackMessages, id)

	smsUserCallbackMessageDto := &models.SmsUserCallbackMessageDto{}
	err = json.Unmarshal([]byte(jsonString), smsUserCallbackMessageDto)
	if nil != err {
		return nil, err
	}

	return smsUserCallbackMessageDto, err
}

func (smsUserCallbackMessageManagement *smsUserCallbackMessageManagement) RemoveById(id string) (bool, error) {
	count, err := global.GetRedisClient().HashDelete(global.RedisKeySmsUserCallbackMessages, id)

	return 0 < count, err
}

func (smsUserCallbackMessageManagement *smsUserCallbackMessageManagement) GetAllIds() ([]string, error) {
	return global.GetRedisClient().HashKeys(global.RedisKeySmsUserCallbackMessages)
}

func (*smsUserCallbackMessageManagement) BuildId(nationCode string, phoneNumber string, createdTime int64) (string) {
	return fmt.Sprintf("%s_%s_%d", nationCode, phoneNumber, createdTime)
}
