package managements

import (
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"encoding/json"
	"time"
	"fmt"
	"errors"
	"rz/middleware/notifycenter/common"
)

var (
	SmsTemplateManagement = smsTemplateManagement{}

	smsTemplateDtos map[int]models.SmsTemplateDto
	lastRefreshTime int64
	refreshDuration int
)

func init() {
	refreshDuration = 60 * 1000
}

type smsTemplateManagement struct {
}

func (smsTemplateManagement *smsTemplateManagement) Set(
	templateId int,
	extend int,
	userCallbackUrls []string,
	pattern string) (error) {
	smsTemplateDto := &models.SmsTemplateDto{
		Id:               templateId,
		Extend:           extend,
		UserCallbackUrls: userCallbackUrls,
		Pattern:          pattern,
	}

	bytes, err := json.Marshal(smsTemplateDto)
	if nil != err {
		return err
	}

	return global.GetRedisClient().HashSet(global.RedisKeySmsTemplates, common.Int32ToString(templateId), string(bytes))
}

func (smsTemplateManagement *smsTemplateManagement) GetByTemplateId(templateId int) (*models.SmsTemplateDto, error) {
	smsTemplateDtos, err := smsTemplateManagement.GetAll()
	if nil != err {
		return nil, err
	}

	smsTemplateDto, ok := smsTemplateDtos[templateId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("map key(%d) is not exist", templateId))
	}

	return &smsTemplateDto, nil
}

func (smsTemplateManagement *smsTemplateManagement) GetByExtend(extend int) (*models.SmsTemplateDto, error) {
	smsTemplateIdExtendMappingDtos, err := smsTemplateManagement.GetAll()
	if nil != err {
		return nil, err
	}

	for _, value := range smsTemplateIdExtendMappingDtos {
		if extend == value.Extend {
			return &value, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("the extend[%d] is not exist", extend))
}

func (smsTemplateManagement *smsTemplateManagement) GetAll() (map[int]models.SmsTemplateDto, error) {
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

func (*smsTemplateManagement) getAll() (map[int]models.SmsTemplateDto, error) {
	jsonStrings, err := global.GetRedisClient().HashGetAll(global.RedisKeySmsTemplates)
	if nil != err {
		return nil, err
	}

	var smsTemplateDtos = make(map[int]models.SmsTemplateDto)
	for _, jsonString := range jsonStrings {
		smsTemplateDto := models.SmsTemplateDto{}
		err = json.Unmarshal([]byte(jsonString), &smsTemplateDto)
		if nil != err {
			continue
		}

		smsTemplateDtos[smsTemplateDto.Id] = smsTemplateDto
	}

	return smsTemplateDtos, nil
}
