package channels

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/core/common"
)

var (
	SmsChannels = map[int]*smsChannelBase{}
)

type smsDoFunc func(smsMessagePo *models.SmsMessagePo, result interface{}) (error)

func ChooseSmsChannel(smsMessagePo *models.SmsMessagePo) (*smsChannelBase, error) {
	err := common.Assert.IsNotNilToError(smsMessagePo, "smsMessagePo")
	if nil != err {
		return nil, err
	}

	smsChannel, ok := SmsChannels[0]
	if !ok {
		return nil, exceptions.FailedChooseSmsChannel().AttachMessage(smsMessagePo.Id)
	}

	return smsChannel, nil
}

type smsChannelBase struct {
	channelBase

	smsDoFunc smsDoFunc
}

func (myself *smsChannelBase) Do(smsMessagePo *models.SmsMessagePo, result interface{}) (error) {
	err := common.Assert.IsNotNilToError(smsMessagePo, "smsMessagePo")
	if nil != err {
		return err
	}
	err = common.Assert.IsNotNilToError(result, "result")
	if nil != err {
		return err
	}

	return myself.smsDoFunc(smsMessagePo, result)
}
