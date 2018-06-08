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
	"rz/middleware/notifycenter/enumerations"
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

	var err error
	SmsMessageConsumer.SendChannel = enumerations.Sms
	SmsMessageConsumer.keySuffix, err = enumerations.SendChannelToString(SmsMessageConsumer.SendChannel)
	common.Assert.IsNilError(err, "")
	SmsMessageConsumer.convertFunc = SmsMessageConsumer.convert
	SmsMessageConsumer.sendFunc = SmsMessageConsumer.Send
}

type smsMessageConsumer struct {
	baseMessageConsumer

	Url               string
	AppKey            string
	AppId             string
	DefaultNationCode string
}

func (smsMessageConsumer *smsMessageConsumer) Send(messageDto interface{}) (error) {
	smsMessageDto := messageDto.(*models.SmsMessageDto)

	var randomNumber = common.Int32ToString(rand.Intn(1024))
	smsMessageRequestExternalDto := smsMessageConsumer.buildSmsMessageRequestExternalDto(smsMessageDto, randomNumber)

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
	smsMessageDto *models.SmsMessageDto,
	randomNumber string) (*external.SmsMessageRequestExternalDto) {
	smsMessageRequestExternalDto := &external.SmsMessageRequestExternalDto{}
	now := time.Now()
	smsMessageRequestExternalDto.TplId = smsMessageDto.TemplateId
	smsMessageRequestExternalDto.Time = now.Unix()
	smsMessageRequestExternalDto.Sig = smsMessageConsumer.buildSignature(smsMessageDto, now, randomNumber)
	smsMessageRequestExternalDto.Tel = smsMessageConsumer.buildPhoneNumberPackExternalDtos(smsMessageDto)
	smsMessageRequestExternalDto.Params = smsMessageDto.Parameters
	if !common.IsStringBlank(smsMessageDto.Content) {
		smsMessageRequestExternalDto.Msg = smsMessageDto.Content
	}
	smsMessageRequestExternalDto.Ext = ""
	smsTemplateDto, err := managements.SmsTemplateManagement.GetByTemplateId(smsMessageDto.TemplateId)
	if nil == err {
		smsMessageRequestExternalDto.Extend = common.Int32ToString(smsTemplateDto.Extend)
	}

	return smsMessageRequestExternalDto
}

func (smsMessageConsumer *smsMessageConsumer) buildSignature(smsMessageDto *models.SmsMessageDto, now time.Time, randomNumber string) (string) {
	var value = fmt.Sprintf(
		"appkey=%s&random=%s&time=%s&mobile=%s",
		smsMessageConsumer.AppKey,
		randomNumber,
		common.Int64ToString(now.Unix()),
		strings.Join(smsMessageDto.Tos, ","))
	var signature = sha256.Sum256([]byte(value))

	return hex.EncodeToString(signature[:])
}

func (smsMessageConsumer *smsMessageConsumer) buildPhoneNumberPackExternalDtos(smsMessageDto *models.SmsMessageDto) ([]external.PhoneNumberPackExternalDto) {
	var phoneNumberPackExternalDtos []external.PhoneNumberPackExternalDto

	var nationCode string
	if common.IsStringBlank(smsMessageDto.NationCode) {
		nationCode = smsMessageConsumer.DefaultNationCode
	}
	for _, phoneNumber := range smsMessageDto.Tos {
		phoneNumberPackExternalDto := external.PhoneNumberPackExternalDto{
			Nationcode: nationCode,
			Mobile:     phoneNumber,
		}
		phoneNumberPackExternalDtos = append(phoneNumberPackExternalDtos, phoneNumberPackExternalDto)
	}

	return phoneNumberPackExternalDtos
}

func (smsMessageConsumer *smsMessageConsumer) convert(jsonString string) (interface{}, *models.BaseMessageDto, error) {
	smsMessageDto := &models.SmsMessageDto{}

	err := json.Unmarshal([]byte(jsonString), smsMessageDto)
	if nil != err {
		return nil, nil, err
	}

	return smsMessageDto, &smsMessageDto.BaseMessageDto, nil
}
