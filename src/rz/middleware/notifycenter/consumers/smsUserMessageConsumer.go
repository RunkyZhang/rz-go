package consumers

import (
	"regexp"
	"time"

	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/enumerations"
	"strings"
	"fmt"
	"errors"
)

var (
	SmsUserMessageConsumer *smsUserMessageConsumer
)

func init() {
	SmsUserMessageConsumer = &smsUserMessageConsumer{
		regularExpressions: make(map[string]*regexp.Regexp),
	}
	SmsUserMessageConsumer.getMessageFunc = SmsUserMessageConsumer.getMessage
	SmsUserMessageConsumer.sendFunc = SmsUserMessageConsumer.Send
	SmsUserMessageConsumer.poToDtoFunc = SmsUserMessageConsumer.poToDto
	SmsUserMessageConsumer.messageManagementBase = &managements.SmsUserMessageManagement.MessageManagementBase
}

type smsUserMessageConsumer struct {
	messageConsumerBase

	regularExpressions map[string]*regexp.Regexp
}

func (myself *smsUserMessageConsumer) Send(messagePo interface{}) (error) {
	smsUserMessagePo, ok := messagePo.(*models.SmsUserMessagePo)
	err := common.Assert.IsTrueToError(ok, "messagePo.(*models.SmsUserMessagePo)")
	if nil != err {
		return err
	}

	smsTemplatePo, err := managements.SmsTemplateManagement.GetByTemplateId(smsUserMessagePo.TemplateId)
	if nil != err {
		return exceptions.TemplateIdNotExist().AttachError(err).AttachMessage(common.Int32ToString(smsUserMessagePo.TemplateId))
	}

	var smsMessagePo *models.SmsMessagePo
	if enumerations.IdentifyingCode == smsTemplatePo.Type {
		smsMessagePo, err = managements.SmsMessageManagement.GetByIdentifyingCode(smsTemplatePo.Id, smsUserMessagePo.Content, time.Now())
		if nil != err {
			return exceptions.InvalidIdentifyingCode().AttachMessage(smsUserMessagePo.Content).AttachError(err)
		}
	} else if enumerations.Pattern == smsTemplatePo.Type {
		regularExpression, ok := myself.regularExpressions[smsTemplatePo.Pattern]
		if !ok {
			regularExpression, err = regexp.Compile(smsTemplatePo.Pattern)
			if nil == err {
				myself.regularExpressions[smsTemplatePo.Pattern] = regularExpression
			} else {
				myself.regularExpressions[smsTemplatePo.Pattern] = nil
			}
		}
		if nil == regularExpression {
			return exceptions.InvalidPattern().AttachMessage(smsTemplatePo.Pattern)
		}
		if !regularExpression.MatchString(smsUserMessagePo.Content) {
			return exceptions.PatternNotMatch().AttachMessage(smsUserMessagePo.Content)
		}
	}

	smsUserCallbackRequestDto := &models.SmsUserCallbackRequestDto{
		Message:     models.SmsMessagePoToDto(smsMessagePo),
		Template:    models.SmsTemplatePoToDto(smsTemplatePo),
		UserMessage: models.SmsUserMessagePoToDto(smsUserMessagePo),
	}

	errorMessages := ""
	urls := strings.Split(smsTemplatePo.UserCallbackUrls, ",")
	for _, url := range urls {
		_, err = httpClient.Post(url, smsUserCallbackRequestDto)
		if nil != err {
			errorMessages += errorMessages + fmt.Sprintf("+++failed to invoke url(%s)", url)
		}
	}

	if "" != errorMessages {
		return errors.New(errorMessages)
	}

	return nil
}

func (myself *smsUserMessageConsumer) getMessage(messageId int, date time.Time) (interface{}, *models.PoBase, *models.CallbackBasePo, error) {
	smsUserMessagePo, err := managements.SmsUserMessageManagement.GetById(messageId, date)
	if nil != err {
		return nil, nil, nil, err
	}

	return smsUserMessagePo, &smsUserMessagePo.PoBase, &smsUserMessagePo.CallbackBasePo, nil
}

func (myself *smsUserMessageConsumer) poToDto(messagePo interface{}) (interface{}) {
	smsUserMessagePo, ok := messagePo.(*models.SmsUserMessagePo)
	if !ok {
		return nil
	}

	return models.SmsUserMessagePoToDto(smsUserMessagePo)
}
