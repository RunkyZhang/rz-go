package provider

import (
	"fmt"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/models/external"
	"rz/core/common"
	"strings"
	"encoding/json"
	"errors"
)

var (
	SmsDahanProvider *smsDahanProvider
)

func init() {
	SmsDahanProvider = &smsDahanProvider{
		Url:      global.GetConfig().SmsDahan.Url,
		Account:  global.GetConfig().SmsDahan.Account,
		Password: global.GetConfig().SmsDahan.Password,
	}
	SmsDahanProvider.smsDoFunc = SmsDahanProvider.do
	SmsDahanProvider.Id = "smsDahanProvider"

	smsProviders[SmsDahanProvider.Id] = &SmsDahanProvider.smsProviderBase
}

type smsDahanProvider struct {
	smsProviderBase

	Url      string
	Account  string
	Password string
}

func (myself *smsDahanProvider) do(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo) (error) {
	smsTemplatePo, err := managements.SmsTemplateManagement.GetByTemplateId(smsMessagePo.TemplateId)
	if nil != err {
		return err
	}

	dahanSmsMessageRequestDto := external.DahanSmsMessageRequestDto{
		Account:  myself.Account,
		Password: myself.Password,
		MsgId:    common.Int64ToString(smsMessagePo.Id),
		Content:  smsMessagePo.Content,
		Sign:     fmt.Sprintf("【%s】", smsTemplatePo.Sign),
		SubCode:  fmt.Sprintf("%d%d", smsTemplatePo.DahanSignCode, smsTemplatePo.Extend),
		SendTime: 201808161630,
		Phones:   myself.buildPhoneNumbers(smsMessagePo.Tos, smsMessagePo.NationCode),
	}

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
			"Result: %s; Desc: %s; FailPhones: %s",
			dahanSmsMessageResponseDto.Result,
			dahanSmsMessageResponseDto.Desc,
			dahanSmsMessageResponseDto.FailPhones)
		return errors.New(message)
	}

	return nil
}

func (myself *smsDahanProvider) buildPhoneNumbers(tos string, nationCode string) (string) {
	value := ""
	if "" != tos {
		phoneNumbers := strings.Split(tos, ",")
		for _, phoneNumber := range phoneNumbers {
			value += fmt.Sprintf("+%s%s,", nationCode, phoneNumber)
		}
	}

	return strings.TrimRight(value, ",")
}
