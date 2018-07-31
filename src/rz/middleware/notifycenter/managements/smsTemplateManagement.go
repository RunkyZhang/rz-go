package managements

import (
	"time"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/exceptions"
)

var (
	SmsTemplateManagement = smsTemplateManagement{
		refreshDuration: 10,
	}

	smsTemplatePos map[int]*models.SmsTemplatePo
)

type smsTemplateManagement struct {
	managementBase

	lastRefreshTime int64
	refreshDuration int
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
		return nil, exceptions.TemplateIdNotExist().AttachError(err).AttachMessage(templateId)
	}

	return smsTemplatePo, nil
}

func (myself *smsTemplateManagement) GetByExtend(extend int) (*models.SmsTemplatePo, error) {
	smsTemplateIdExtendMappingPos, err := myself.GetAll()
	if nil != err {
		return nil, err
	}

	for _, value := range smsTemplateIdExtendMappingPos {
		if extend == value.Extend {
			return value, nil
		}
	}

	return nil, exceptions.ExtendNotExist().AttachError(err).AttachMessage(extend)
}

func (myself *smsTemplateManagement) GetAll() (map[int]*models.SmsTemplatePo, error) {
	var err error
	if nil == smsTemplatePos {
		smsTemplatePos, err = myself.getAll()
		return smsTemplatePos, err
	}

	if int64(myself.refreshDuration) <= time.Now().Unix()-myself.lastRefreshTime {
		go func() {
			smsTemplatePos, err = myself.getAll()
			if nil == err {
				myself.lastRefreshTime = time.Now().Unix()
			} else {
				common.GetLogging().Error(err, "Failed to get all [SmsTemplatePo]")
			}
		}()
	}

	return smsTemplatePos, nil
}

func (*smsTemplateManagement) ExistExtend(extend int) (bool, error) {
	count, err := repositories.SmsTemplateRepository.CountByExtend(extend)

	return 0 < count, err
}

func (*smsTemplateManagement) getAll() (map[int]*models.SmsTemplatePo, error) {
	smsTemplatePos, err := repositories.SmsTemplateRepository.SelectAll()
	if nil != err {
		return nil, err
	}

	var keyValues = make(map[int]*models.SmsTemplatePo)
	for _, smsTemplatePo := range smsTemplatePos {
		keyValues[smsTemplatePo.Id] = smsTemplatePo
	}

	return keyValues, nil
}
