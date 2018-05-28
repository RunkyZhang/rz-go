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
	"rz/middleware/notifycenter/enumerations"
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

	SmsMessageConsumer.SendChannel = enumerations.Sms
	SmsMessageConsumer.consumeFunc = SmsMessageConsumer.consume
	SmsMessageConsumer.handleErrorFunc = SmsMessageConsumer.handleError
}

type smsMessageConsumer struct {
	baseMessageConsumer

	Url               string
	AppKey            string
	AppId             string
	DefaultNationCode int
}

func (smsMessageConsumer *smsMessageConsumer) Send(smsMessageDto *models.SmsMessageDto) error {
	var randomNumber = strconv.Itoa(rand.Intn(1024))
	var smsMessageRequestExternalDto = &external.SmsMessageRequestExternalDto{}
	smsMessageRequestExternalDto.Time = time.Now().Unix()
	var sig = fmt.Sprintf(
		"appkey=%s&random=%s&time=%s&mobile=%s",
		smsMessageConsumer.AppKey,
		randomNumber,
		strconv.FormatInt(smsMessageRequestExternalDto.Time, 10),
		strings.Join(smsMessageDto.Tos, ","))
	var sigsha256 = sha256.Sum256([]byte(sig))
	smsMessageRequestExternalDto.Sig = hex.EncodeToString(sigsha256[:])
	smsMessageRequestExternalDto.Tel = []external.PhoneNumberExternalDto{}
	for _, phoneNumber := range smsMessageDto.Tos {
		phoneNumberExternalDto := external.PhoneNumberExternalDto{
			Nationcode: strconv.Itoa(smsMessageConsumer.DefaultNationCode),
			Mobile:     phoneNumber,
		}
		smsMessageRequestExternalDto.Tel = append(smsMessageRequestExternalDto.Tel, phoneNumberExternalDto)
	}
	smsMessageRequestExternalDto.Type = "0"
	smsMessageRequestExternalDto.Msg = fmt.Sprintf("[应用告警]%s", smsMessageDto.Content)
	smsMessageRequestExternalDto.Ext = ""
	smsMessageRequestExternalDto.Extend = ""

	bytes, err := httplib.Post(smsMessageConsumer.Url+"?sdkappid="+smsMessageConsumer.AppId+"&random="+randomNumber, smsMessageRequestExternalDto)
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

func (smsMessageConsumer *smsMessageConsumer) consume(jsonString string) (interface{}, error) {
	smsMessageDto := &models.SmsMessageDto{}

	err := json.Unmarshal([]byte(jsonString), smsMessageDto)
	if nil != err {
		return smsMessageDto, nil
	}

	return smsMessageDto, smsMessageConsumer.Send(smsMessageDto)
}

func (*smsMessageConsumer) handleError(messageDto interface{}, err error) (error) {
	smsMessageDto := messageDto.(*models.SmsMessageDto)

	sendChannel, err := enumerations.SendChannelToString(SmsMessageConsumer.SendChannel)
	if nil != err {
		return err
	}
	messageState, err := enumerations.MessageStateToString(enumerations.Error)
	if nil != err {
		messageState = "Unknown"
	}

	smsMessageDto.States = smsMessageDto.States + ">" + messageState
	smsMessageDto.Exception = err.Error()

	bytes, err := json.Marshal(smsMessageDto)
	if nil != err {
		return err
	}

	return global.GetRedisClient().HashSet(global.RedisKeyMessageValues+sendChannel, smsMessageDto.Id, string(bytes))
}
