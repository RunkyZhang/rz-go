package consumers

import (
	"time"
	"rz/middleware/notifycenter/managements"
	"fmt"
	"rz/middleware/notifycenter/models"
	"regexp"
)

var (
	regularExpressions map[string]*regexp.Regexp
)

type smsUserCallbackConsumer struct {
}

func (smsUserCallbackConsumer *smsUserCallbackConsumer) Start(duration time.Duration) {
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			smsUserCallbackConsumer.start()
			timer.Reset(duration)
		}
	}
}

func (smsUserCallbackConsumer *smsUserCallbackConsumer) start() {
	smsUserCallbackMessageIds, err := managements.SmsUserMessageManagement.GetAllIds()
	if nil != err {
		fmt.Printf("failed to get [smsUserCallbackMessage] ids. error: %s", err.Error())
		return
	}
	if nil == smsUserCallbackMessageIds {
		return
	}

	for _, smsUserCallbackMessageId := range smsUserCallbackMessageIds {
		smsUserCallbackMessageDto, err := managements.SmsUserMessageManagement.GetById(smsUserCallbackMessageId)
		if nil != err {
			fmt.Printf("failed to get [UserCallbackMessage](%s) value. error: %s", smsUserCallbackMessageId, err.Error())
			_, err := managements.SmsUserMessageManagement.RemoveById(smsUserCallbackMessageId)
			if nil != err {
				fmt.Printf("failed to remove [UserCallbackMessage](%s). error: %s", smsUserCallbackMessageId, err.Error())
			}

			continue
		}

		ok, err := managements.SmsUserMessageManagement.RemoveById(smsUserCallbackMessageId)
		if nil != err || !ok {
			fmt.Printf("failed to remove [UserCallbackMessage](%s). error: %s", smsUserCallbackMessageId, err.Error())

			continue
		}

		smsUserCallbackDto, err := managements.SmsUserCallbackManagement.GetById(smsUserCallbackMessageDto.UserCallbackId)
		if nil != err {
			fmt.Printf("failed to remove [UserCallback](%s). error: %s", smsUserCallbackMessageDto.UserCallbackId, err.Error())

			continue
		}
		smsTemplateDto, err := managements.SmsTemplateManagement.GetByTemplateId(smsUserCallbackDto.TemplateId)
		if nil != err {
			fmt.Printf("failed to remove [SmsTemplate](%d). error: %s", smsUserCallbackDto.TemplateId, err.Error())

			continue
		}

		if nil == smsUserCallbackDto.UserCallbackMessages {
			smsUserCallbackDto.UserCallbackMessages = make(map[string]models.SmsUserMessageDto)
		}
		//if smsTemplateDto.Disable {
		//	smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
		//	smsUserCallbackMessageDto.Finished = true
		//	smsUserCallbackMessageDto.ErrorMessage = "Disable Template"
		//	smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
		//	managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)
		//
		//	continue
		//}

		regularExpression, ok := regularExpressions[smsTemplateDto.Pattern]
		if !ok {
			regularExpression, err = regexp.Compile(smsTemplateDto.Pattern)
			if nil == err {
				regularExpressions[smsTemplateDto.Pattern] = regularExpression
			} else {
				regularExpressions[smsTemplateDto.Pattern] = nil
			}
		}
		if nil == regularExpression {
			smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
			smsUserCallbackMessageDto.Finished = true
			smsUserCallbackMessageDto.ErrorMessage = "Invalid Pattern"
			smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
			managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)

			continue
		}
		if !regularExpression.MatchString(smsUserCallbackMessageDto.Content) {
			smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
			smsUserCallbackMessageDto.Finished = true
			smsUserCallbackMessageDto.ErrorMessage = "Not Match"
			smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
			managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)

			continue
		}

		if nil == smsTemplateDto.UserCallbackUrls {
			smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
			smsUserCallbackMessageDto.Finished = true
			smsUserCallbackMessageDto.ErrorMessage = "Nil Callback Urls"
			smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
			managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)

			continue
		}

		for _, userCallbackUrl := range smsTemplateDto.UserCallbackUrls {
			fmt.Printf("send %s", userCallbackUrl)
		}

		smsUserCallbackMessageDto.FinishedTime = time.Now().Unix()
		smsUserCallbackMessageDto.Finished = true
		smsUserCallbackMessageDto.ErrorMessage = "Sent"
		smsUserCallbackDto.UserCallbackMessages[smsUserCallbackMessageDto.Id] = *smsUserCallbackMessageDto
		managements.SmsUserCallbackManagement.Set(smsUserCallbackDto)
	}
}
