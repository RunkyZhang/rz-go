package services

import (
	"time"
	"fmt"
	"strings"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/managements"
	"rz/middleware/notifycenter/enumerations"
	"rz/middleware/notifycenter/exceptions"
	"rz/core/common"
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
	systemAliasPermissionPo, err := managements.SystemAliasPermissionManagement.GetById(smsMessageDto.SystemAlias)
	if nil != err || 0 == systemAliasPermissionPo.SmsPermission {
		return 0, exceptions.NotSendSmsPermission().AttachError(err).AttachMessage(smsMessageDto.SystemAlias)
	}
	smsTemplatePo, err := managements.SmsTemplateManagement.GetByTemplateId(smsMessageDto.TemplateId)
	if nil != err {
		return 0, exceptions.InvalidTemplateId().AttachError(err).AttachMessage(smsMessageDto.TemplateId)
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
			enumerations.Initial,
			enumerations.Error,
			exceptions.FailedEnqueueMessageId().AttachError(err).AttachMessage(smsMessagePo.Id).Error(),
			"",
			&finished,
			&now,
			smsMessagePo.CreatedTime.Year())

		return 0, exceptions.FailedEnqueueMessageId().AttachError(err).AttachMessage(smsMessagePo.Id)
	}

	return smsMessagePo.Id, err
}

func (myself *smsMessageService) QueryByIds(queryMessagesByIdsRequestDto *models.QueryMessagesByIdsRequestDto) ([]*models.SmsMessageDto, error) {
	err := common.Assert.IsTrueToError(nil != queryMessagesByIdsRequestDto, "nil != queryMessagesByIdsRequestDto")
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
