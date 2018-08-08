package channels

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/core/common"
)

var (
	SmsUserChannels = map[int]*smsUserChannelBase{}
)

type smsUserDoFunc func(smsUserMessagePo *models.SmsUserMessagePo) (error)

func ChooseSmsUserChannel(smsUserMessagePo *models.SmsUserMessagePo) (*smsUserChannelBase, error) {
	err := common.Assert.IsNotNilToError(smsUserMessagePo, "smsUserMessagePo")
	if nil != err {
		return nil, err
	}

	smsUserChannel, ok := SmsUserChannels[0]
	if !ok {
		return nil, exceptions.FailedChooseSmsUserChannel().AttachMessage(smsUserMessagePo.Id)
	}

	return smsUserChannel, nil
}

type smsUserChannelBase struct {
	channelBase

	smsUserDoFunc smsUserDoFunc
}

func (myself *smsUserChannelBase) Do(smsUserMessagePo *models.SmsUserMessagePo) (error) {
	err := common.Assert.IsNotNilToError(smsUserMessagePo, "smsUserMessagePo")
	if nil != err {
		return err
	}

	return myself.smsUserDoFunc(smsUserMessagePo)
}
