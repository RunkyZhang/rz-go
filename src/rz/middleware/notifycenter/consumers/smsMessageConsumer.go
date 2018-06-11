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

	"git.zhaogangren.com/cloud/cloud.base.utils-go.sdk/httplib"

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

	SmsMessageConsumer.convertFunc = SmsMessageConsumer.convert
	SmsMessageConsumer.sendFunc = SmsMessageConsumer.Send
	SmsMessageConsumer.messageManagementBase = &managements.SmsMessageManagement.MessageManagementBase
}

type smsMessageConsumer struct {
	messageConsumerBase

	Url               string
	AppKey            string
	AppId             string
	DefaultNationCode string
}

func (smsMessageConsumer *smsMessageConsumer) Send(messagePo interface{}) (error) {
	smsMessagePo := messagePo.(*models.SmsMessagePo)

	var randomNumber = common.Int32ToString(rand.Intn(1024))
	smsMessageRequestExternalDto := smsMessageConsumer.buildSmsMessageRequestExternalDto(smsMessagePo, randomNumber)

	url := fmt.Sprintf("%s?sdkappid=%s&random=%s", smsMessageConsumer.Url, smsMessageConsumer.AppId, randomNumber)
	bytes, err := httplib.Post(url, smsMessageRequestExternalDto)
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

func (smsMessageConsumer *smsMessageConsumer) buildSmsMessageRequestExternalDto(
	smsMessagePo *models.SmsMessagePo,
	randomNumber string) (*external.SmsMessageRequestExternalDto) {
	smsMessageRequestExternalDto := &external.SmsMessageRequestExternalDto{}
	now := time.Now()
	smsMessageRequestExternalDto.TplId = smsMessagePo.TemplateId
	smsMessageRequestExternalDto.Time = now.Unix()
	smsMessageRequestExternalDto.Sig = smsMessageConsumer.buildSignature(smsMessagePo, now, randomNumber)
	smsMessageRequestExternalDto.Tel = smsMessageConsumer.buildPhoneNumberPackExternalDtos(smsMessagePo)
	smsMessageRequestExternalDto.Params = strings.Split(smsMessagePo.Parameters, ",")
	if !common.IsStringBlank(smsMessagePo.Content) {
		smsMessageRequestExternalDto.Msg = smsMessagePo.Content
	}
	smsMessageRequestExternalDto.Ext = ""
	smsTemplateDto, err := managements.SmsTemplateManagement.GetByTemplateId(smsMessagePo.TemplateId)
	if nil == err {
		smsMessageRequestExternalDto.Extend = common.Int32ToString(smsTemplateDto.Extend)
	}

	return smsMessageRequestExternalDto
}

func (smsMessageConsumer *smsMessageConsumer) buildSignature(smsMessagePo *models.SmsMessagePo, now time.Time, randomNumber string) (string) {
	var value = fmt.Sprintf(
		"appkey=%s&random=%s&time=%s&mobile=%s",
		smsMessageConsumer.AppKey,
		randomNumber,
		common.Int64ToString(now.Unix()),
		smsMessagePo.Tos)
	var signature = sha256.Sum256([]byte(value))

	return hex.EncodeToString(signature[:])
}

func (smsMessageConsumer *smsMessageConsumer) buildPhoneNumberPackExternalDtos(smsMessagePo *models.SmsMessagePo) ([]external.PhoneNumberPackExternalDto) {
	var phoneNumberPackExternalDtos []external.PhoneNumberPackExternalDto

	var nationCode string
	if common.IsStringBlank(smsMessagePo.NationCode) {
		nationCode = smsMessageConsumer.DefaultNationCode
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

func (smsMessageConsumer *smsMessageConsumer) convert(messageId int) (interface{}, *models.MessageBasePo, error) {
	smsMessageDto, err := managements.SmsMessageManagement.GetById(messageId)
	if nil != err {
		return nil, nil, err
	}

	return smsMessageDto, &smsMessageDto.MessageBasePo, nil
}
