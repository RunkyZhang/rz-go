package provider

import (
	"fmt"
	"encoding/json"
	"strings"
	"errors"

	"rz/core/common"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/managements"
)

var (
	SmsDahanProvider *smsDahanProvider
)

func init() {
	SmsDahanProvider = &smsDahanProvider{}
	SmsDahanProvider.smsDoFunc = SmsDahanProvider.do
	SmsDahanProvider.Id = "smsDahanProvider"
	smsProviderPo, err := managements.SmsProviderManagement.GetById(SmsDahanProvider.Id)
	common.Assert.IsNilErrorToPanic(err, "Failed to get [SmsProviderPo]")
	SmsDahanProvider.Url = smsProviderPo.Url1
	keyValues := make(map[string]string)
	err = json.Unmarshal([]byte(smsProviderPo.PassportJson), &keyValues)
	var ok bool
	SmsDahanProvider.Account, ok = keyValues["account"]
	if !ok {
		common.Assert.IsTrueToPanic(ok, "map has not [account]")
	}
	SmsDahanProvider.Password, ok = keyValues["password"]
	if !ok {
		common.Assert.IsTrueToPanic(ok, "map has not [password]")
	}

	smsProviders[SmsDahanProvider.Id] = &SmsDahanProvider.smsProviderBase
}

type smsDahanProvider struct {
	smsProviderBase

	Url      string
	Account  string
	Password string
}

func (myself *smsDahanProvider) do(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo) (error) {
	dahanSmsMessageRequestDto := myself.buildDahanSmsMessageRequestDto(smsMessagePo, smsTemplatePo)

	bytes, err := global.HttpClient.Post(myself.Url, dahanSmsMessageRequestDto)
	if nil != err {
		return err
	}

	dahanSmsMessageResponseDto := &external.DahanSmsMessageResponseDto{}
	err = json.Unmarshal(bytes, dahanSmsMessageResponseDto)
	if nil != err {
		return err
	}

	if "0" != dahanSmsMessageResponseDto.Result {
		message := fmt.Sprintf(
			"Result: %s; MsgId: %s; Desc: %s; FailPhones: %s",
			dahanSmsMessageResponseDto.Result,
			dahanSmsMessageResponseDto.MsgId,
			dahanSmsMessageResponseDto.Desc,
			dahanSmsMessageResponseDto.FailPhones)
		return errors.New(message)
	}

	return nil
}

func (myself *smsDahanProvider) buildDahanSmsMessageRequestDto(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo) (*external.DahanSmsMessageRequestDto) {
	dahanSmsMessageRequestDto := &external.DahanSmsMessageRequestDto{}
	dahanSmsMessageRequestDto.Password = myself.Password
	dahanSmsMessageRequestDto.Account = myself.Account
	dahanSmsMessageRequestDto.Content = smsMessagePo.Content
	dahanSmsMessageRequestDto.MsgId = common.Int64ToString(smsMessagePo.Id)
	dahanSmsMessageRequestDto.Sign = fmt.Sprintf("【%s】", smsTemplatePo.Sign)
	dahanSmsMessageRequestDto.SendTime = 201808161630
	dahanSmsMessageRequestDto.SubCode = fmt.Sprintf("%d%d", smsTemplatePo.DahanSignCode, smsTemplatePo.Extend)
	dahanSmsMessageRequestDto.Phones = myself.buildPhoneNumbers(smsMessagePo)

	return dahanSmsMessageRequestDto
}

func (myself *smsDahanProvider) buildPhoneNumbers(smsMessagePo *models.SmsMessagePo) (string) {
	nationCode := smsMessagePo.NationCode
	if "" == smsMessagePo.NationCode {
		nationCode = global.DefaultNationCode
	}

	value := ""
	if "" != smsMessagePo.Tos {
		phoneNumbers := strings.Split(smsMessagePo.Tos, ",")
		for _, phoneNumber := range phoneNumbers {
			value += fmt.Sprintf("+%s%s,", nationCode, phoneNumber)
		}
		value = value[0:(len(value) - 1)]
	}

	return value
}
