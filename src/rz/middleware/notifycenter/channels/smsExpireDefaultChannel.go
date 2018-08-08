package channels

import (
	"time"
	"strings"
	"errors"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/global"
	"rz/core/common"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/exceptions"
)

var (
	SmsExpireDefaultChannel *smsExpireDefaultChannel
)

func init() {
	SmsExpireDefaultChannel = &smsExpireDefaultChannel{}
	SmsExpireDefaultChannel.smsExpireDoFunc = SmsExpireDefaultChannel.do
	SmsExpireDefaultChannel.Id = 0

	SmsExpireChannels[SmsExpireDefaultChannel.Id] = &SmsExpireDefaultChannel.smsExpireChannelBase
}

type smsExpireDefaultChannel struct {
	smsExpireChannelBase
}

func (myself *smsExpireDefaultChannel) do(smsMessagePo *models.SmsMessagePo) (error) {
	err := common.Assert.IsNotNilToError(smsMessagePo, "smsMessagePo")
	if nil != err {
		return err
	}

	if "" == smsMessagePo.ExpireCallbackUrls {
		return nil
	}

	smsUserMessagePos, err := managements.SmsUserMessageManagement.Query(
		smsMessagePo.Id, "", "", "", 0, time.Now().Year())
	if nil != err {
		return err
	}
	if 0 != len(smsUserMessagePos) {
		return nil
	}

	messageExpireCallbackRequestDto := &models.MessageExpireCallbackRequestDto{
		Message: models.SmsMessagePoToDto(smsMessagePo),
	}
	errorMessage := ""
	urls := strings.Split(smsMessagePo.ExpireCallbackUrls, ",")
	for _, url := range urls {
		_, err = global.HttpClient.Post(url, messageExpireCallbackRequestDto)
		if nil != err {
			errorMessage += "[" + exceptions.FailedRequestHttp().AttachError(err).AttachMessage(url).Error() + "]"
		}
	}

	if "" != errorMessage {
		return errors.New(errorMessage)
	}

	return nil
}
