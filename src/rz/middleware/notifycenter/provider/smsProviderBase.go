package provider

import (
	"rz/core/common"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"errors"
	"math/rand"
	"fmt"
)

var (
	smsProviders = map[string]*smsProviderBase{}
)

type validSmsProvider struct {
	SmsProvider   *smsProviderBase
	SmsProviderPo *models.SmsProviderPo
	Weighted      int
	Start         int
	End           int
	Rate          int
}

func ChooseSmsProvider(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo, excludedIds []string) (*smsProviderBase, error) {
	err := common.Assert.IsTrueToError(nil != smsMessagePo, "nil != smsMessagePo")
	if nil != err {
		return nil, err
	}
	err = common.Assert.IsTrueToError(nil != smsTemplatePo, "nil != smsTemplatePo")
	if nil != err {
		return nil, err
	}

	totalWeighted := 0
	var validSmsProviders []*validSmsProvider
	for _, smsProvider := range smsProviders {
		excluded := false
		if nil != excludedIds && 0 != len(excludedIds) {
			for _, excludedId := range excludedIds {
				if smsProvider.Id == excludedId {
					excluded = true
					break
				}
			}
		}
		if excluded {
			continue
		}

		smsProviderPo, err := managements.SmsProviderManagement.GetById(smsProvider.Id)
		if nil != err {
			continue
		}

		if validateSmsProvider(smsTemplatePo, smsProviderPo) {
			validSmsProvider := &validSmsProvider{
				SmsProvider:   smsProvider,
				SmsProviderPo: smsProviderPo,
			}
			validSmsProvider.Weighted = validSmsProvider.SmsProviderPo.Weighted
			if 0 >= validSmsProvider.Weighted {
				validSmsProvider.Weighted = 1
			}
			validSmsProviders = append(validSmsProviders, validSmsProvider)
			totalWeighted += validSmsProvider.Weighted
		}
	}

	if 0 == len(validSmsProviders) {
		return nil, errors.New("there is no any valid [SmsProvider]")
	} else if 1 == len(validSmsProviders) {
		return validSmsProviders[0].SmsProvider, nil
	}

	position := 0
	for _, validSmsProvider := range validSmsProviders {
		rate := int((float32(validSmsProvider.Weighted) / float32(totalWeighted)) * 100)
		if 0 == rate {
			rate = 1
		}
		validSmsProvider.Rate = rate
		validSmsProvider.Start = position
		position += rate
		validSmsProvider.End = position
	}

	smsProvider := validSmsProviders[0].SmsProvider
	randomNumber := rand.Intn(position)
	for _, validSmsProvider := range validSmsProviders {
		if randomNumber >= validSmsProvider.Start && randomNumber <= validSmsProvider.End {
			smsProvider = validSmsProvider.SmsProvider
			break
		}
	}

	return smsProvider, nil
}

func validateSmsProvider(smsTemplatePo *models.SmsTemplatePo, smsProviderPo *models.SmsProviderPo) (bool) {
	if smsProviderPo.Disable {
		return false
	}
	if smsTemplatePo.ContentType != (smsTemplatePo.ContentType & smsProviderPo.ContentTypes) {
		return false
	}

	return true
}

func testWeighted() {
	var validSmsProviders []*validSmsProvider
	validSmsProvider1 := &validSmsProvider{
		Weighted: 2,
	}
	validSmsProviders = append(validSmsProviders, validSmsProvider1)
	validSmsProvider1 = &validSmsProvider{
		Weighted: 3,
	}
	validSmsProviders = append(validSmsProviders, validSmsProvider1)
	validSmsProvider1 = &validSmsProvider{
		Weighted: 555,
	}
	validSmsProviders = append(validSmsProviders, validSmsProvider1)

	totalWeighted := 0
	for _, validSmsProvider := range validSmsProviders {
		totalWeighted += validSmsProvider.Weighted
	}
	position := 0
	for _, validSmsProvider := range validSmsProviders {
		rate := int((float32(validSmsProvider.Weighted) / float32(totalWeighted)) * 100)
		if 0 == rate {
			rate = 1
		}
		validSmsProvider.Rate = rate
		validSmsProvider.Start = position
		position += rate
		validSmsProvider.End = position
	}

	keyValues := make(map[string]int)
	for i := 0; i < 100000; i++ {
		randomNumber := rand.Intn(position)
		for _, validSmsProvider := range validSmsProviders {
			if randomNumber >= validSmsProvider.Start && randomNumber <= validSmsProvider.End {
				key := common.Int32ToString(validSmsProvider.Weighted) + "--" + common.Int32ToString(validSmsProvider.Rate)
				keyValues[key] = keyValues[key] + 1
				break
			}
		}
	}
	fmt.Println(keyValues)
}

type smsDoFunc func(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo) (error)

type smsProviderBase struct {
	providerBase

	smsDoFunc smsDoFunc
}

func (myself *smsProviderBase) Do(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo) (error) {
	err := common.Assert.IsTrueToError(nil != smsMessagePo, "nil != smsMessagePo")
	if nil != err {
		return err
	}
	err = common.Assert.IsTrueToError(nil != smsTemplatePo, "nil != smsTemplatePo")
	if nil != err {
		return err
	}

	return myself.smsDoFunc(smsMessagePo, smsTemplatePo)
}
