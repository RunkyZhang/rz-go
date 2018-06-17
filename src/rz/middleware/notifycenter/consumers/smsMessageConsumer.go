package consumers

import (
	"math/rand"
	"time"
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/managements"
)

var (
	SmsMessageConsumer *smsMessageConsumer
)

func init() {
	SmsMessageConsumer = &smsMessageConsumer{
		Url:               global.Config.Sms.Url,
		AppId:             global.Config.Sms.AppId,
		AppKey:            global.Config.Sms.AppKey,
		DefaultNationCode: global.Config.Sms.DefaultNationCode,
	}
	SmsMessageConsumer.getMessageFunc = SmsMessageConsumer.getMessage
	SmsMessageConsumer.sendFunc = SmsMessageConsumer.Send
	SmsMessageConsumer.poToDtoFunc = SmsMessageConsumer.poToDto
	SmsMessageConsumer.messageManagementBase = &managements.SmsMessageManagement.MessageManagementBase
	SmsMessageConsumer.httpClient = common.NewHttpClient()
}

type smsMessageConsumer struct {
	messageConsumerBase

	Url               string
	AppKey            string
	AppId             string
	DefaultNationCode string
	httpClient        *common.HttpClient
}

func (myself *smsMessageConsumer) Send(messagePo interface{}) (error) {
	smsMessagePo := messagePo.(*models.SmsMessagePo)

	var randomNumber = common.Int32ToString(rand.Intn(1024))
	smsMessageRequestExternalDto := myself.buildSmsMessageRequestExternalDto(smsMessagePo, randomNumber)

	url := fmt.Sprintf("%s?sdkappid=%s&random=%s", myself.Url, myself.AppId, randomNumber)
	bytes, err := myself.httpClient.Post(url, smsMessageRequestExternalDto)
	if nil != err {
		return err
	}

	smsMessageResponseExternalDto := &external.SmsMessageResponseExternalDto{}
	err = json.Unmarshal(bytes, smsMessageResponseExternalDto)
	if nil != err {
		return err
	}
	if 0 != smsMessageResponseExternalDto.Result {
		return errors.New(smsMessageResponseExternalDto.Errmsg)
	}

	return nil
}

func (myself *smsMessageConsumer) buildSmsMessageRequestExternalDto(
	smsMessagePo *models.SmsMessagePo,
	randomNumber string) (*external.SmsMessageRequestExternalDto) {
	smsMessageRequestExternalDto := &external.SmsMessageRequestExternalDto{}
	now := time.Now()
	smsMessageRequestExternalDto.TplId = smsMessagePo.TemplateId
	smsMessageRequestExternalDto.Time = now.Unix()
	smsMessageRequestExternalDto.Sig = myself.buildSignature(smsMessagePo, now, randomNumber)
	smsMessageRequestExternalDto.Tel = myself.buildPhoneNumberPackExternalDtos(smsMessagePo)
	smsMessageRequestExternalDto.Params = strings.Split(smsMessagePo.Parameters, ",")
	if !common.IsStringBlank(smsMessagePo.Content) {
		smsMessageRequestExternalDto.Msg = smsMessagePo.Content
	}
	smsMessageRequestExternalDto.Ext = ""
	smsTemplatePo, err := managements.SmsTemplateManagement.GetByTemplateId(smsMessagePo.TemplateId)
	if nil == err {
		smsMessageRequestExternalDto.Extend = common.Int32ToString(smsTemplatePo.Extend)
	}

	return smsMessageRequestExternalDto
}

func (myself *smsMessageConsumer) buildSignature(smsMessagePo *models.SmsMessagePo, now time.Time, randomNumber string) (string) {
	var value = fmt.Sprintf(
		"appkey=%s&random=%s&time=%s&mobile=%s",
		myself.AppKey,
		randomNumber,
		common.Int64ToString(now.Unix()),
		smsMessagePo.Tos)
	var signature = sha256.Sum256([]byte(value))

	return hex.EncodeToString(signature[:])
}

func (myself *smsMessageConsumer) buildPhoneNumberPackExternalDtos(smsMessagePo *models.SmsMessagePo) ([]external.PhoneNumberPackExternalDto) {
	var phoneNumberPackExternalDtos []external.PhoneNumberPackExternalDto

	nationCode := smsMessagePo.NationCode
	if "" == smsMessagePo.NationCode {
		nationCode = myself.DefaultNationCode
	}

	phoneNumbers := strings.Split(smsMessagePo.Tos, ",")
	for _, phoneNumber := range phoneNumbers {
		phoneNumberPackExternalDto := external.PhoneNumberPackExternalDto{
			Nationcode: nationCode,
			Mobile:     phoneNumber,
		}
		phoneNumberPackExternalDtos = append(phoneNumberPackExternalDtos, phoneNumberPackExternalDto)
	}

	return phoneNumberPackExternalDtos
}

func (myself *smsMessageConsumer) getMessage(messageId int, date time.Time) (interface{}, *models.PoBase, *models.CallbackBasePo, error) {
	smsMessagePo, err := managements.SmsMessageManagement.GetById(messageId, date)
	if nil != err {
		return nil, nil, nil, err
	}

	return smsMessagePo, &smsMessagePo.PoBase, &smsMessagePo.CallbackBasePo, nil
}

func (myself *smsMessageConsumer) poToDto(messagePo interface{}) (interface{}) {
	smsMessagePo := messagePo.(*models.SmsMessagePo)

	return models.SmsMessagePoToDto(smsMessagePo)
}
