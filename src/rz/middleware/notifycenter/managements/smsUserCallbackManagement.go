package managements

import (
	"rz/middleware/notifycenter/global"
	"fmt"
	"rz/middleware/notifycenter/models"
	"encoding/json"
	"rz/middleware/notifycenter/common"
)

var (
	SmsUserCallbackManagement = smsUserCallbackManagement{}
)

type smsUserCallbackManagement struct {
}

func (smsUserCallbackManagement *smsUserCallbackManagement) Get(nationCode string, phoneNumber string, extend int) (*models.SmsUserCallbackDto, error) {
	smsTemplate, err := SmsTemplateManagement.GetByExtend(extend)
	if nil != err {
		return nil, err
	}

	id := smsUserCallbackManagement.buildId(nationCode, phoneNumber, smsTemplate.Id)

	return smsUserCallbackManagement.GetById(id)
}

func (smsUserCallbackManagement *smsUserCallbackManagement) GetById(id string) (*models.SmsUserCallbackDto, error) {
	jsonString, err := global.GetRedisClient().HashGet(global.RedisKeySmsUserCallbcaks, id)
	if nil != err {
		return nil, err
	}

	smsUserCallbackDto := &models.SmsUserCallbackDto{}
	err = json.Unmarshal([]byte(jsonString), smsUserCallbackDto)
	if nil != err {
		return nil, err
	}

	return smsUserCallbackDto, nil
}

func (smsUserCallbackManagement *smsUserCallbackManagement) Set(smsUserCallbackDto *models.SmsUserCallbackDto) (error) {
	if nil == smsUserCallbackDto {
		return common.Assert.NewNilParameterError("smsUserCallbackDto")
	}

	bytes, err := json.Marshal(smsUserCallbackDto)
	if nil != err {
		return err
	}

	return global.GetRedisClient().HashSet(global.RedisKeySmsUserCallbcaks, smsUserCallbackDto.Id, string(bytes))
}

func (*smsUserCallbackManagement) buildId(nationCode string, phoneNumber string, templateId int) (string) {
	return fmt.Sprintf("%s_%s_%d", nationCode, phoneNumber, templateId)
}
