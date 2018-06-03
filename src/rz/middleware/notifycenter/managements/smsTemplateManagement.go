package managements

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"encoding/json"
	"time"
	"fmt"
	"errors"
)

var (
	SmsTemplateManagement = smsTemplateManagement{}

	smsTemplateDtos map[int]*models.SmsTemplateDto
	lastRefreshTime int64
	refreshDuration int
)

func init() {
	refreshDuration = 60 * 1000
}

type smsTemplateManagement struct {
}

func (smsTemplateManagement *smsTemplateManagement) GetByTemplateId(templateId int) (*models.SmsTemplateDto, error) {
	smsTemplateDtos, err := smsTemplateManagement.GetAll()
	if nil != err {
		return nil, err
	}

	smsTemplateDto, err := smsTemplateDtos[templateId]
	if nil != err {
		return nil, err
	}

	return smsTemplateDto, nil
}

func (smsTemplateManagement *smsTemplateManagement) GetByExtend(extend int) (*models.SmsTemplateDto, error) {
	smsTemplateIdExtendMappingDtos, err := smsTemplateManagement.GetAll()
	if nil != err {
		return nil, err
	}

	for _, value := range smsTemplateIdExtendMappingDtos {
		if extend == value.Extend {
			return value, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("the extend[%d] is not exist", extend))
}

func (smsTemplateManagement *smsTemplateManagement) GetAll() (map[int]*models.SmsTemplateDto, error) {
	var err error
	if nil == smsTemplateDtos {
		smsTemplateDtos, err = smsTemplateManagement.getAll()
		return smsTemplateDtos, err
	}

	if int64(refreshDuration) <= time.Now().Unix()-lastRefreshTime {
		go func() {
			smsTemplateDtos, err = smsTemplateManagement.getAll()
			if nil == err {
				lastRefreshTime = time.Now().Unix()
			} else {
				fmt.Println("failed to get all [smsTemplateIdExtendMappingDto], error:", err)
			}
		}()
	}

	return smsTemplateDtos, nil
}

func (*smsTemplateManagement) getAll() (map[int]*models.SmsTemplateDto, error) {
	jsonStrings, err := global.GetRedisClient().HashGetAll(global.RedisKeySmsTemplates)
	if nil != err {
		return nil, err
	}

	var smsTemplateDtos = make(map[int]*models.SmsTemplateDto)
	for _, jsonString := range jsonStrings {
		smsTemplateDto := &models.SmsTemplateDto{}
		err = json.Unmarshal([]byte(jsonString), smsTemplateDto)
		if nil != err {
			continue
		}

		smsTemplateDtos[smsTemplateDto.Id] = smsTemplateDto
	}

	return smsTemplateDtos, nil
}
