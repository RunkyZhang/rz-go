package services

import (
	"time"
	"fmt"
	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/common"
	"strings"
)

var (
	SmsMessageService = smsMessageService{}
)

func init() {
	SmsMessageService.messageManagementBase = &managements.SmsMessageManagement.MessageManagementBase
}

type smsMessageService struct {
	MessageServiceBase
}

func (myself *smsMessageService) Send(smsMessageDto *models.SmsMessageDto) (int64, error) {
	err := VerifySmsMessageDto(smsMessageDto)
	if nil != err {
		return 0, err
	}
	smsTemplatePo, err := managements.SmsTemplateManagement.GetByTemplateId(smsMessageDto.TemplateId)
	if nil != err {
		return 0, err
	}
	count := len(strings.Split(smsTemplatePo.Context, "%s")) - 1
	if count != len(smsMessageDto.Parameters) {
		return 0, exceptions.InvalidSmsParameterCount().AttachMessage(fmt.Sprintf("%d != %d", count, len(smsMessageDto.Parameters)))
	}

	smsMessagePo := models.SmsMessageDtoToPo(smsMessageDto)
	smsMessagePo.Content = myself.calculateContext(smsTemplatePo.Context, smsMessageDto.Parameters)
	smsMessagePo.CreatedTime = time.Now()
	smsMessagePo.Id, err = managements.SmsMessageManagement.GenerateId(smsMessagePo.CreatedTime.Year())
	if nil != err {
		return 0, exceptions.FailedGenerateMessageId().AttachError(err)
	}

	if "" != smsMessagePo.ExpireCallbackUrls {
		err = managements.SmsMessageManagement.EnqueueExpireIds(smsMessagePo.Id, smsMessagePo.ExpireTime.Unix())
		if nil != err {
			return 0, exceptions.FailedEnqueueExpireMessageId().AttachError(err)
		}
	}

	err = managements.SmsMessageManagement.Add(smsMessagePo)
	if nil != err {
		return 0, err
	}

	err = managements.SmsMessageManagement.EnqueueIds(smsMessagePo.Id, smsMessagePo.ScheduleTime.Unix())
	if nil != err {
		now := time.Now()
		finished := true
		managements.ModifyMessageFlowAsync(
			myself.messageManagementBase,
			smsMessagePo.Id,
			enumerations.Error,
			exceptions.FailedEnqueueMessageId().AttachError(err).AttachMessage(smsMessagePo.Id).Error(),
			&finished,
			&now,
			smsMessagePo.CreatedTime.Year())

		return 0, err
	}

	return smsMessagePo.Id, err
}

func (myself *smsMessageService) QueryByIds(queryMessagesByIdsRequestDto *models.QueryMessagesByIdsRequestDto) ([]*models.SmsMessageDto, error) {
	err := common.Assert.IsNotNilToError(queryMessagesByIdsRequestDto, "queryMessagesByIdsRequestDto")
	if nil != err {
		return nil, err
	}

	idGroups := map[int][]int64{}
	for _, value := range queryMessagesByIdsRequestDto.Ids {
		if 4 > len(value) {
			common.GetLogging().Warn(err, exceptions.InvalidMessageId().AttachMessage(value).Error())
			continue
		}
		year, err := common.StringToInt32(value[0:4])
		if nil != err {
			common.GetLogging().Warn(err, exceptions.InvalidMessageId().AttachMessage(value).Error())
			continue
		}
		id, err := common.StringToInt64(value)
		if nil != err {
			common.GetLogging().Warn(err, exceptions.InvalidMessageId().AttachMessage(value).Error())
			continue
		}

		_, ok := idGroups[year]
		if !ok {
			idGroups[year] = []int64{}
		}
		idGroup, _ := idGroups[year]
		idGroups[year] = append(idGroup, id)
		fmt.Println(idGroup)
	}

	var smsMessagePos []*models.SmsMessagePo
	for year, idGroup := range idGroups {
		partSmsMessagePos, err := managements.SmsMessageManagement.GetByIds(idGroup, year)
		if nil != err {
			common.GetLogging().Warn(err, exceptions.DatabaseError().AttachMessage(idGroup).Error())
			continue
		}

		smsMessagePos = append(smsMessagePos, partSmsMessagePos...)
	}

	return models.SmsMessagePosToDtos(smsMessagePos), nil
}

func (myself *smsMessageService) calculateContext(context string, parameters []string) (string) {
	var args []interface{}
	for _, parameter := range parameters {
		args = append(args, parameter)
	}

	return fmt.Sprintf(context, args...)
}
