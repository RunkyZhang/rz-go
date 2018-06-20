package managements

import (
	"rz/middleware/notifycenter/models"
	"time"
	"fmt"
	"errors"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/common"
)

var (
	SmsTemplateManagement = smsTemplateManagement{}

	smsTemplatePos  map[int]models.SmsTemplatePo
	lastRefreshTime int64
	refreshDuration int
)

func init() {
	refreshDuration = 60 * 1000
}

type smsTemplateManagement struct {
	managementBase
}

func (myself *smsTemplateManagement) Add(smsTemplatePo *models.SmsTemplatePo) (error) {
	err := common.Assert.IsNotNilToError(smsTemplatePo, "smsTemplatePo")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(smsTemplatePo.PoBase, "smsTemplatePo.PoBase")
	if nil != err {
		return err
	}

	myself.setPoBase(&smsTemplatePo.PoBase)

	return repositories.SmsTemplateRepository.Insert(smsTemplatePo)
}

func (myself *smsTemplateManagement) GetByTemplateId(templateId int) (*models.SmsTemplatePo, error) {
	smsTemplatePos, err := myself.GetAll()
	if nil != err {
		return nil, err
	}

	smsTemplatePo, ok := smsTemplatePos[templateId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("map key(%d) is not exist", templateId))
	}

	return &smsTemplatePo, nil
}

func (myself *smsTemplateManagement) GetByExtend(extend int) (*models.SmsTemplatePo, error) {
	smsTemplateIdExtendMappingPos, err := myself.GetAll()
	if nil != err {
		return nil, err
	}

	for _, value := range smsTemplateIdExtendMappingPos {
		if extend == value.Extend {
			return &value, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("the extend[%d] is not exist", extend))
}

func (myself *smsTemplateManagement) GetAll() (map[int]models.SmsTemplatePo, error) {
	var err error
	if nil == smsTemplatePos {
		smsTemplatePos, err = myself.getAll()
		return smsTemplatePos, err
	}

	if int64(refreshDuration) <= time.Now().Unix()-lastRefreshTime {
		go func() {
			smsTemplatePos, err = myself.getAll()
			if nil == err {
				lastRefreshTime = time.Now().Unix()
			} else {
				fmt.Printf("failed to get all [smsTemplatePo], error: %s\n", err.Error())
			}
		}()
	}

	return smsTemplatePos, nil
}

func (*smsTemplateManagement) getAll() (map[int]models.SmsTemplatePo, error) {
	smsTemplatePos, err := repositories.SmsTemplateRepository.SelectAll()
	if nil != err {
		return nil, err
	}

	var keyValues = make(map[int]models.SmsTemplatePo)
	for _, smsTemplatePo := range smsTemplatePos {
		keyValues[smsTemplatePo.Id] = smsTemplatePo
	}

	return keyValues, nil
}
