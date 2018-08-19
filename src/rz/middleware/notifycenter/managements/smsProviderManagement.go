package managements

import (
	"time"

	"rz/core/common"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/repositories"
	"rz/middleware/notifycenter/exceptions"
)

var (
	SmsProviderManagement = smsProviderManagement{
		refreshDuration: 10,
	}

	smsProviderPos map[string]*models.SmsProviderPo
)

type smsProviderManagement struct {
	managementBase

	lastRefreshTime int64
	refreshDuration int
}

func (myself *smsProviderManagement) Add(smsProviderPo *models.SmsProviderPo) (error) {
	err := common.Assert.IsTrueToError(nil != smsProviderPo, "nil != smsProviderPo")
	if nil != err {
		return err
	}

	myself.setPoBase(&smsProviderPo.PoBase)

	return repositories.SmsProviderRepository.Insert(smsProviderPo)
}

func (myself *smsProviderManagement) GetById(id string) (*models.SmsProviderPo, error) {
	smsProviderPos, err := myself.GetAll()
	if nil != err {
		return nil, err
	}

	smsProviderPo, ok := smsProviderPos[id]
	if !ok {
		return nil, exceptions.TemplateIdNotExist().AttachError(err).AttachMessage(id)
	}

	return smsProviderPo, nil
}

func (myself *smsProviderManagement) GetAll() (map[string]*models.SmsProviderPo, error) {
	var err error
	if nil == smsProviderPos {
		smsProviderPos, err = myself.getAll()
		return smsProviderPos, err
	}

	if int64(myself.refreshDuration) <= time.Now().Unix()-myself.lastRefreshTime {
		go func() {
			smsProviderPos, err = myself.getAll()
			if nil == err {
				myself.lastRefreshTime = time.Now().Unix()
			} else {
				common.GetLogging().Error(err, "Failed to get all [SmsProviderPo]")
			}
		}()
	}

	return smsProviderPos, nil
}

func (*smsProviderManagement) getAll() (map[string]*models.SmsProviderPo, error) {
	smsProviderPos, err := repositories.SmsProviderRepository.SelectAll()
	if nil != err {
		return nil, err
	}

	var keyValues = make(map[string]*models.SmsProviderPo)
	for _, smsProviderPo := range smsProviderPos {
		keyValues[smsProviderPo.Id] = smsProviderPo
	}

	return keyValues, nil
}
