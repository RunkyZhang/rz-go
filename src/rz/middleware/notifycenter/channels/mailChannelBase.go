package channels

import (
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
)

var (
	MailChannels = map[int]*mailChannelBase{}
)

type mailDoFunc func(mailMessagePo *models.MailMessagePo) (error)

func ChooseMailChannel(mailMessagePo *models.MailMessagePo) (*mailChannelBase, error) {
	err := common.Assert.IsNotNilToError(mailMessagePo, "mailMessagePo")
	if nil != err {
		return nil, err
	}

	mailChannel, ok := MailChannels[0]
	if !ok {
		return nil, exceptions.FailedChooseMailChannel().AttachMessage(mailMessagePo.Id)
	}

	return mailChannel, nil
}

type mailChannelBase struct {
	channelBase

	mailDoFunc mailDoFunc
}

func (myself *mailChannelBase) Do(mailMessagePo *models.MailMessagePo) (error) {
	err := common.Assert.IsNotNilToError(mailMessagePo, "mailMessagePo")
	if nil != err {
		return err
	}

	return myself.mailDoFunc(mailMessagePo)
}
