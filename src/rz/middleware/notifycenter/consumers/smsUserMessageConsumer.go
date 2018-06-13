package consumers

import (
	"regexp"
	"time"

	"git.zhaogangren.com/cloud/cloud.base.utils-go.sdk/httplib"

	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/enumerations"
)

var (
	SmsUserMessageConsumer *smsUserMessageConsumer
)

func init() {
	SmsUserMessageConsumer = &smsUserMessageConsumer{
		regularExpressions: make(map[string]*regexp.Regexp),
	}
	SmsUserMessageConsumer.convertFunc = SmsUserMessageConsumer.convert
	SmsUserMessageConsumer.sendFunc = SmsUserMessageConsumer.Send
	SmsUserMessageConsumer.messageManagementBase = &managements.SmsUserMessageManagement.MessageManagementBase
}

type smsUserMessageConsumer struct {
	messageConsumerBase

	regularExpressions map[string]*regexp.Regexp
}

func (myself *smsUserMessageConsumer) Send(messagePo interface{}) (error) {
	smsUserMessagePo := messagePo.(*models.SmsUserMessagePo)

	smsTemplatePo, err := managements.SmsTemplateManagement.GetByTemplateId(smsUserMessagePo.TemplateId)
	if nil != err {
		return exceptions.TemplateIdNotExist().AttachError(err).AttachMessage(common.Int32ToString(smsUserMessagePo.TemplateId))
	}

	var smsMessagePo *models.SmsMessagePo
	if enumerations.IdentifyingCode == smsTemplatePo.Type {
		smsMessagePo, err = managements.SmsMessageManagement.GetByIdentifyingCode(smsTemplatePo.Id, smsUserMessagePo.Content, time.Now())
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
			return exceptions.InvalidPattern().AttachMessage(smsUserMessagePo.Content)
		}
	}

	smsUserCallbackRequestDto := &models.SmsUserCallbackRequestDto{
		Message:     models.SmsMessagePoToDto(smsMessagePo),
		Template:    models.SmsTemplatePoToDto(smsTemplatePo),
		UserMessage: models.SmsUserMessagePoToDto(smsUserMessagePo),
	}
	for _, userCallbackUrl := range smsTemplatePo.UserCallbackUrls {
		httplib.Post(userCallbackUrl, smsUserCallbackRequestDto)
	}

	return nil
}

func (myself *smsUserMessageConsumer) convert(messageId int, date time.Time) (interface{}, *models.PoBase, *models.CallbackBasePo, error) {
	smsUserMessagePo, err := managements.SmsUserMessageManagement.GetById(messageId, date)
	if nil != err {
		return nil, nil, nil, err
	}

	return smsUserMessagePo, &smsUserMessagePo.PoBase, &smsUserMessagePo.CallbackBasePo, nil
}

//func (myself *smsUserCallbackConsumer) Start(duration time.Duration) {
//	timer := time.NewTimer(duration)
//
//	for {
//		select {
//		case <-timer.C:
//			myself.start()
//			timer.Reset(duration)
//		}
//	}
//}
//
//func (myself *smsUserCallbackConsumer) start() {
//	now := time.Now()
//	userMessageIds, err := managements.SmsUserMessageManagement.DequeueMessageIds(now)
//	if nil != err || nil == userMessageIds {
//		fmt.Println("failed to get user message ids. error: ", err)
//		return
//	}
//
//	for _, userMessageId := range userMessageIds {
//		smsUserCallbackMessageDto, err := managements.SmsUserMessageManagement.GetById(userMessageId, now)
//		if nil != err {
//			fmt.Printf("failed to get [UserCallbackMessage](%s) value. error: %s", userMessageId, err.Error())
//			_, err := managements.SmsUserMessageManagement.RemoveById(userMessageId)
//			if nil != err {
//				fmt.Printf("failed to remove [UserCallbackMessage](%s). error: %s", userMessageId, err.Error())
//			}
//
//			continue
//		}
//
//		ok, err := managements.SmsUserMessageManagement.RemoveById(userMessageId)
//		if nil != err || !ok {
//			fmt.Printf("failed to remove [UserCallbackMessage](%s). error: %s", userMessageId, err.Error())
//
//			continue
//		}
//
//		smsUserCallbackDto, err := managements.SmsUserCallbackManagement.GetById(smsUserCallbackMessageDto.UserCallbackId)
//		if nil != err {
//			fmt.Printf("failed to remove [UserCallback](%s). error: %s", smsUserCallbackMessageDto.UserCallbackId, err.Error())
//
//			continue
//		}
//		smsTemplateDto, err := managements.SmsTemplateManagement.GetByTemplateId(smsUserCallbackDto.TemplateId)
//		if nil != err {
//			fmt.Printf("failed to remove [SmsTemplate](%d). error: %s", smsUserCallbackDto.TemplateId, err.Error())
//
//			continue
//		}
//
//		if nil == smsUserCallbackDto.UserCallbackMessages {
//			smsUserCallbackDto.UserCallbackMessages = make(map[string]models.SmsUserMessageDto)
//		}
//		//if smsTemplateDto.Disable {
//		//	smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
//		//	smsUserCallbackMessageDto.Finished = true
//		//	smsUserCallbackMessageDto.ErrorMessage = "Disable Template"
//		//	smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
//		//	managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)
//		//
//		//	continue
//		//}
//
//		regularExpression, ok := regularExpressions[smsTemplateDto.Pattern]
//		if !ok {
//			regularExpression, err = regexp.Compile(smsTemplateDto.Pattern)
//			if nil == err {
//				regularExpressions[smsTemplateDto.Pattern] = regularExpression
//			} else {
//				regularExpressions[smsTemplateDto.Pattern] = nil
//			}
//		}
//		if nil == regularExpression {
//			smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
//			smsUserCallbackMessageDto.Finished = true
//			smsUserCallbackMessageDto.ErrorMessage = "Invalid Pattern"
//			smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
//			managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)
//
//			continue
//		}
//		if !regularExpression.MatchString(smsUserCallbackMessageDto.Content) {
//			smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
//			smsUserCallbackMessageDto.Finished = true
//			smsUserCallbackMessageDto.ErrorMessage = "Not Match"
//			smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
//			managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)
//
//			continue
//		}
//
//		if nil == smsTemplateDto.UserCallbackUrls {
//			smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
//			smsUserCallbackMessageDto.Finished = true
//			smsUserCallbackMessageDto.ErrorMessage = "Nil Callback Urls"
//			smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
//			managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)
//
//			continue
//		}
//
//		for _, userCallbackUrl := range smsTemplateDto.UserCallbackUrls {
//			fmt.Printf("send %s", userCallbackUrl)
//		}
//
//		smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
//		smsUserCallbackMessageDto.Finished = true
//		smsUserCallbackMessageDto.ErrorMessage = "Sent"
//		smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
//		managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)
//	}
//}
