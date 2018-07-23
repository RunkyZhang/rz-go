package channels

import (
	"math/rand"
	"fmt"
	"encoding/json"
	"time"
	"strings"
	"crypto/sha256"
	"encoding/hex"

	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/models/external"
	"rz/middleware/notifycenter/managements"
)

var (
	SmsTencentChannel *smsTencentChannel
)

func init() {
	SmsTencentChannel = &smsTencentChannel{
		Url:               global.GetConfig().Sms.Url,
		AppId:             global.GetConfig().Sms.AppId,
		AppKey:            global.GetConfig().Sms.AppKey,
		DefaultNationCode: global.GetConfig().Sms.DefaultNationCode,
	}
	SmsTencentChannel.smsDoFunc = SmsTencentChannel.do
	SmsTencentChannel.Id = 0

	SmsChannels[SmsTencentChannel.Id] = &SmsTencentChannel.smsChannelBase
}

type smsTencentChannel struct {
	smsChannelBase

	Url               string
	AppKey            string
	AppId             string
	DefaultNationCode string
}

func (myself *smsTencentChannel) do(smsMessagePo *models.SmsMessagePo, result interface{}) (error) {
	smsTemplatePo, err := managements.SmsTemplateManagement.GetByTemplateId(smsMessagePo.TemplateId)
	if nil != err {
		return err
	}

	var randomNumber = common.Int32ToString(rand.Intn(1024))
	smsMessageRequestExternalDto := myself.buildSmsMessageRequestExternalDto(smsMessagePo, smsTemplatePo, randomNumber)

	url := fmt.Sprintf("%s?sdkappid=%s&random=%s", myself.Url, myself.AppId, randomNumber)
	bytes, err := global.HttpClient.Post(url, smsMessageRequestExternalDto)
	if nil != err {
		return err
	}

	err = json.Unmarshal(bytes, result)
	if nil != err {
		return err
	}

	return nil
}

func (myself *smsTencentChannel) buildSmsMessageRequestExternalDto(
	smsMessagePo *models.SmsMessagePo,
	smsTemplatePo *models.SmsTemplatePo,
	randomNumber string) (*external.SmsMessageRequestExternalDto) {
	smsMessageRequestExternalDto := &external.SmsMessageRequestExternalDto{}
	now := time.Now()
	smsMessageRequestExternalDto.TplId = smsTemplatePo.TencentTemplateId
	smsMessageRequestExternalDto.Time = now.Unix()
	smsMessageRequestExternalDto.Sig = myself.buildSignature(smsMessagePo, now, randomNumber)
	smsMessageRequestExternalDto.Tel = myself.buildPhoneNumberPackExternalDtos(smsMessagePo)
	if "" != smsMessagePo.Parameters {
		smsMessageRequestExternalDto.Params = strings.Split(smsMessagePo.Parameters, ",")
	}
	//if !common.IsStringBlank(smsMessagePo.Content) {
	//	smsMessageRequestExternalDto.Msg = smsMessagePo.Content
	//}
	smsMessageRequestExternalDto.Ext = ""
	smsMessageRequestExternalDto.Extend = common.Int32ToString(smsTemplatePo.Extend)

	return smsMessageRequestExternalDto
}

func (myself *smsTencentChannel) buildSignature(smsMessagePo *models.SmsMessagePo, now time.Time, randomNumber string) (string) {
	var value = fmt.Sprintf(
		"appkey=%s&random=%s&time=%s&mobile=%s",
		myself.AppKey,
		randomNumber,
		common.Int64ToString(now.Unix()),
		smsMessagePo.Tos)
	var signature = sha256.Sum256([]byte(value))

	return hex.EncodeToString(signature[:])
}

func (myself *smsTencentChannel) buildPhoneNumberPackExternalDtos(smsMessagePo *models.SmsMessagePo) ([]external.PhoneNumberPackExternalDto) {
	var phoneNumberPackExternalDtos []external.PhoneNumberPackExternalDto

	nationCode := smsMessagePo.NationCode
	if "" == smsMessagePo.NationCode {
		nationCode = myself.DefaultNationCode
	}

	if "" != smsMessagePo.Tos {
		phoneNumbers := strings.Split(smsMessagePo.Tos, ",")
		for _, phoneNumber := range phoneNumbers {
			phoneNumberPackExternalDto := external.PhoneNumberPackExternalDto{
				Nationcode: nationCode,
				Mobile:     phoneNumber,
			}
			phoneNumberPackExternalDtos = append(phoneNumberPackExternalDtos, phoneNumberPackExternalDto)
		}
	}

	return phoneNumberPackExternalDtos
}
