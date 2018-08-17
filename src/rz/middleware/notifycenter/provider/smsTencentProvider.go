package provider

import (
	"math/rand"
	"fmt"
	"encoding/json"
	"time"
	"strings"
	"crypto/sha256"
	"encoding/hex"

	"rz/core/common"
	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/models/external"
	"errors"
)

var (
	SmsTencentProvider *smsTencentProvider
)

func init() {
	SmsTencentProvider = &smsTencentProvider{
		Url:               global.GetConfig().SmsTencent.Url,
		AppId:             global.GetConfig().SmsTencent.AppId,
		AppKey:            global.GetConfig().SmsTencent.AppKey,
		DefaultNationCode: global.GetConfig().SmsTencent.DefaultNationCode,
	}
	SmsTencentProvider.smsDoFunc = SmsTencentProvider.do
	SmsTencentProvider.Id = "smsTencentProvider"

	smsProviders[SmsTencentProvider.Id] = &SmsTencentProvider.smsProviderBase
}

type smsTencentProvider struct {
	smsProviderBase

	Url               string
	AppKey            string
	AppId             string
	DefaultNationCode string
}

func (myself *smsTencentProvider) do(smsMessagePo *models.SmsMessagePo, smsTemplatePo *models.SmsTemplatePo) (error) {
	var randomNumber = common.Int32ToString(rand.Intn(1024))
	tencentSmsMessageRequestDto := myself.buildTencentSmsMessageRequestDto(smsMessagePo, smsTemplatePo, randomNumber)

	url := fmt.Sprintf("%s?sdkappid=%s&random=%s", myself.Url, myself.AppId, randomNumber)
	bytes, err := global.HttpClient.Post(url, tencentSmsMessageRequestDto)
	if nil != err {
		return err
	}

	tencentSmsMessageResponseDto := &external.TencentSmsMessageResponseDto{}
	err = json.Unmarshal(bytes, tencentSmsMessageResponseDto)
	if nil != err {
		return err
	}

	if 0 != tencentSmsMessageResponseDto.ErrorCode {
		message := fmt.Sprintf(
			"ErrorInfo: %s; ActionStatus: %s; ErrorCode: %d",
			tencentSmsMessageResponseDto.ErrorInfo,
			tencentSmsMessageResponseDto.ActionStatus,
			tencentSmsMessageResponseDto.ErrorCode)
		return errors.New(message)
	}
	if 0 != tencentSmsMessageResponseDto.Result {
		message := fmt.Sprintf(
			"Errmsg: %s; Result: %d",
			tencentSmsMessageResponseDto.Errmsg,
			tencentSmsMessageResponseDto.Result)
		return errors.New(message)
	}

	return nil
}

func (myself *smsTencentProvider) buildTencentSmsMessageRequestDto(
	smsMessagePo *models.SmsMessagePo,
	smsTemplatePo *models.SmsTemplatePo,
	randomNumber string) (*external.TencentSmsMessageRequestDto) {
	tencentSmsMessageRequestDto := &external.TencentSmsMessageRequestDto{}
	now := time.Now()
	tencentSmsMessageRequestDto.TplId = smsTemplatePo.TencentTemplateId
	tencentSmsMessageRequestDto.Time = now.Unix()
	tencentSmsMessageRequestDto.Sig = myself.buildSignature(smsMessagePo, now, randomNumber)
	tencentSmsMessageRequestDto.Tel = myself.buildTencentPhoneNumberPackDtos(smsMessagePo)
	if "" != smsMessagePo.Parameters {
		tencentSmsMessageRequestDto.Params = strings.Split(smsMessagePo.Parameters, ",")
	}
	//if !common.IsStringBlank(smsMessagePo.Content) {
	//	tencentSmsMessageRequestDto.Msg = smsMessagePo.Content
	//}
	tencentSmsMessageRequestDto.Ext = ""
	tencentSmsMessageRequestDto.Extend = common.Int32ToString(smsTemplatePo.Extend)

	return tencentSmsMessageRequestDto
}

func (myself *smsTencentProvider) buildSignature(smsMessagePo *models.SmsMessagePo, now time.Time, randomNumber string) (string) {
	var value = fmt.Sprintf(
		"appkey=%s&random=%s&time=%s&mobile=%s",
		myself.AppKey,
		randomNumber,
		common.Int64ToString(now.Unix()),
		smsMessagePo.Tos)
	var signature = sha256.Sum256([]byte(value))

	return hex.EncodeToString(signature[:])
}

func (myself *smsTencentProvider) buildTencentPhoneNumberPackDtos(smsMessagePo *models.SmsMessagePo) ([]external.TencentPhoneNumberPackDto) {
	var tencentPhoneNumberPackDtos []external.TencentPhoneNumberPackDto

	nationCode := smsMessagePo.NationCode
	if "" == smsMessagePo.NationCode {
		nationCode = myself.DefaultNationCode
	}

	if "" != smsMessagePo.Tos {
		phoneNumbers := strings.Split(smsMessagePo.Tos, ",")
		for _, phoneNumber := range phoneNumbers {
			tencentPhoneNumberPackDto := external.TencentPhoneNumberPackDto{
				Nationcode: nationCode,
				Mobile:     phoneNumber,
			}
			tencentPhoneNumberPackDtos = append(tencentPhoneNumberPackDtos, tencentPhoneNumberPackDto)
		}
	}

	return tencentPhoneNumberPackDtos
}
