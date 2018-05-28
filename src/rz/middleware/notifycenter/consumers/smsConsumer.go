package consumers

import (
	"strconv"
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
)

var (
	SmsConsumer *smsConsumer
)

func init() {
	SmsConsumer = &smsConsumer{
		Url:               global.Config.Sms.Url,
		AppId:             global.Config.Sms.AppId,
		AppKey:            global.Config.Sms.AppKey,
		DefaultNationCode: global.Config.Sms.DefaultNationCode,
	}
}

type smsConsumer struct {
	Url               string
	AppKey            string
	AppId             string
	DefaultNationCode int
}

func (smsConsumer *smsConsumer) Send(smsMessageDto models.SmsMessageDto) error {
	var randomNumber = strconv.Itoa(rand.Intn(1024))
	var smsMessageRequestExternalDto = &external.SmsMessageRequestExternalDto{}
	smsMessageRequestExternalDto.Time = time.Now().Unix()
	var sig = fmt.Sprintf(
		"appkey=%s&random=%s&time=%s&mobile=%s",
		smsConsumer.AppKey,
		randomNumber,
		strconv.FormatInt(smsMessageRequestExternalDto.Time, 10),
		strings.Join(smsMessageDto.Tos, ","))
	var sigsha256 = sha256.Sum256([]byte(sig))
	smsMessageRequestExternalDto.Sig = hex.EncodeToString(sigsha256[:])
	smsMessageRequestExternalDto.Tel = []external.PhoneNumberExternalDto{}
	for _, phoneNumber := range smsMessageDto.Tos {
		phoneNumberExternalDto := external.PhoneNumberExternalDto{
			Nationcode: string(smsConsumer.DefaultNationCode),
			Mobile:     phoneNumber,
		}
		smsMessageRequestExternalDto.Tel = append(smsMessageRequestExternalDto.Tel, phoneNumberExternalDto)
	}
	smsMessageRequestExternalDto.Type = "0"
	smsMessageRequestExternalDto.Msg = fmt.Sprintf("[应用告警]%s", smsMessageDto.Content)
	smsMessageRequestExternalDto.Ext = ""
	smsMessageRequestExternalDto.Extend = ""

	bytes, err := httplib.Post(smsConsumer.Url+"?sdkappid="+smsConsumer.AppId+"&random="+randomNumber, smsMessageRequestExternalDto)
	if nil != err {
		return err
	}
	var smsMessageResponseExternalDto external.SmsMessageResponseExternalDto
	err = json.Unmarshal(bytes, &smsMessageResponseExternalDto)
	if nil != err {
		return err
	}
	if 0 == smsMessageResponseExternalDto.Result {
		return nil
	}

	return errors.New(smsMessageResponseExternalDto.Errmsg)
}
