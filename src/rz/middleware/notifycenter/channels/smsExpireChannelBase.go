package channels

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
)

var (
	SmsExpireChannels = map[int]*smsExpireChannelBase{}
)

type smsExpireDoFunc func(smsMessagePo *models.SmsMessagePo) (error)

func ChooseSmsExpireChannel(smsMessagePo *models.SmsMessagePo) (*smsExpireChannelBase, error) {
	err := common.Assert.IsNotNilToError(smsMessagePo, "smsMessagePo")
	if nil != err {
		return nil, err
	}

	smsExpireChannel, ok := SmsExpireChannels[0]
	if !ok {
		return nil, exceptions.FailedChooseSmsExpireChannel().AttachMessage(smsMessagePo.Id)
	}

	return smsExpireChannel, nil
}

type smsExpireChannelBase struct {
	channelBase

	smsExpireDoFunc smsExpireDoFunc
}

func (myself *smsExpireChannelBase) Do(smsMessagePo *models.SmsMessagePo) (error) {
	err := common.Assert.IsNotNilToError(smsMessagePo, "smsMessagePo")
	if nil != err {
		return err
	}

	return myself.smsExpireDoFunc(smsMessagePo)
}
