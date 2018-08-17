package provider

import (
	"regexp"
	"time"
	"strings"
	"errors"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/global"
)

var (
	SmsUserDefaultProvider *smsUserDefaultProvider
)

func init() {
	SmsUserDefaultProvider = &smsUserDefaultProvider{
		regularExpressions: make(map[string]*regexp.Regexp),
	}
	SmsUserDefaultProvider.smsUserDoFunc = SmsUserDefaultProvider.do
	SmsUserDefaultProvider.Id = "smsUserDefaultProvider"

	smsUserProviders[SmsUserDefaultProvider.Id] = &SmsUserDefaultProvider.smsUserProviderBase
}

type smsUserDefaultProvider struct {
	smsUserProviderBase

	regularExpressions map[string]*regexp.Regexp
}

func (myself *smsUserDefaultProvider) do(smsUserMessagePo *models.SmsUserMessagePo) (error) {
	smsTemplatePo, err := managements.SmsTemplateManagement.GetByTemplateId(smsUserMessagePo.TemplateId)
	if nil != err {
		return exceptions.TemplateIdNotExist().AttachError(err).AttachMessage(smsUserMessagePo.TemplateId)
	}

	var smsMessagePo *models.SmsMessagePo
	if enumerations.IdentifyingCode == smsTemplatePo.Type {
		smsMessagePo, err = managements.SmsMessageManagement.GetByIdentifyingCode(smsTemplatePo.Id, smsUserMessagePo.Content, time.Now().Year())
		if nil != err {
			return exceptions.InvalidIdentifyingCode().AttachMessage(smsUserMessagePo.Content).AttachError(err)
		}
		_, err = managements.SmsUserMessageManagement.ModifySmsMessageId(smsUserMessagePo.Id, smsMessagePo.Id)
		if nil != err {
			return exceptions.FailedModifySmsMessageId().AttachMessage(smsUserMessagePo.Id).AttachError(err)
		}

		if smsMessagePo.NationCode != smsUserMessagePo.NationCode || !strings.Contains(smsMessagePo.Tos, smsUserMessagePo.PhoneNumber) {
			return exceptions.FailedMatchPhoneNumber().AttachMessage(smsUserMessagePo.PhoneNumber)
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

	errorMessage := ""
	if "" != smsTemplatePo.UserCallbackUrls {
		urls := strings.Split(smsTemplatePo.UserCallbackUrls, ",")
		for _, url := range urls {
			_, err = global.HttpClient.Post(url, smsUserCallbackRequestDto)
			if nil != err {
				errorMessage += "[" + exceptions.FailedRequestHttp().AttachError(err).AttachMessage(url).Error() + "]"
			}
		}
	}

	if "" != errorMessage {
		return errors.New(errorMessage)
	}

	return nil
}
