package exceptions

import "rz/core/common"

type newBusinessErrorFunc func() (*common.BusinessError)

var (
	Ok                           = func() (*common.BusinessError) { return common.NewBusinessError("Ok", 0) }
	InternalServerError          = func() (*common.BusinessError) { return common.NewBusinessError("Internal server error", 1) }
	DtoNull                      = func() (*common.BusinessError) { return common.NewBusinessError("[Dto] is null", 30001) }
	InvalidDtoType               = func() (*common.BusinessError) { return common.NewBusinessError("Invalid [Dto] type", 30002) }
	TosEmpty                     = func() (*common.BusinessError) { return common.NewBusinessError("[Tos] is empty", 30003) }
	SubjectBlank                 = func() (*common.BusinessError) { return common.NewBusinessError("The subject is blank", 30004) }
	InvalidSendChannel           = func() (*common.BusinessError) { return common.NewBusinessError("Invalid send channel number", 30005) }
	InvalidMessageState          = func() (*common.BusinessError) { return common.NewBusinessError("Invalid message state", 30006) }
	MessageExpire                = func() (*common.BusinessError) { return common.NewBusinessError("Message expire", 30007) }
	MessageBodyMissed            = func() (*common.BusinessError) { return common.NewBusinessError("Message body is missed", 30008) }
	FailedAddSmsUserMessage      = func() (*common.BusinessError) { return common.NewBusinessError("Failed to add sms user message", 30009) }
	FailedEnqueueMessageId       = func() (*common.BusinessError) { return common.NewBusinessError("Failed to enqueue message id", 30010) }
	InvalidExtend                = func() (*common.BusinessError) { return common.NewBusinessError("Invalid extend", 30011) }
	TemplateIdNotExist           = func() (*common.BusinessError) { return common.NewBusinessError("Template id is not exist", 30012) }
	InvalidPattern               = func() (*common.BusinessError) { return common.NewBusinessError("Invalid pattern", 30013) }
	PatternNotMatch              = func() (*common.BusinessError) { return common.NewBusinessError("Pattern is not match string", 30014) }
	FailedInvokeController       = func() (*common.BusinessError) { return common.NewBusinessError("Failed invoke controller", 30015) }
	InvalidSystemAlias           = func() (*common.BusinessError) { return common.NewBusinessError("Invalid [SystemAlias]", 30016) }
	InvalidMessageExpireTime     = func() (*common.BusinessError) { return common.NewBusinessError("Invalid message expire time", 30017) }
	InvalidIdentifyingCode       = func() (*common.BusinessError) { return common.NewBusinessError("Invalid identifying code", 30018) }
	MessageSystemAliasMotMatch   = func() (*common.BusinessError) { return common.NewBusinessError("Message [systemAlias] is not match", 30019) }
	MessageDisable               = func() (*common.BusinessError) { return common.NewBusinessError("Message is disabled", 30020) }
	ExtendExist                  = func() (*common.BusinessError) { return common.NewBusinessError("Extend is exist", 30021) }
	MessageNotInitialState       = func() (*common.BusinessError) { return common.NewBusinessError("Message is not initial state", 30022) }
	FailedChooseSmsChannel       = func() (*common.BusinessError) { return common.NewBusinessError("Failed to choose Sms channel", 30023) }
	ExtendNotExist               = func() (*common.BusinessError) { return common.NewBusinessError("Extend is not exist", 30024) }
	FailedGenerateMessageId      = func() (*common.BusinessError) { return common.NewBusinessError("Failed to generate message id", 30025) }
	InvalidMessageId             = func() (*common.BusinessError) { return common.NewBusinessError("Invalid message id", 30026) }
	DatabaseError                = func() (*common.BusinessError) { return common.NewBusinessError("Database error", 30027) }
	NullQueryParameter           = func() (*common.BusinessError) { return common.NewBusinessError("Null query parameter", 30028) }
	FailedRequestHttp            = func() (*common.BusinessError) { return common.NewBusinessError("Failed request http", 30029) }
	FailedChooseSmsUserChannel   = func() (*common.BusinessError) { return common.NewBusinessError("Failed to choose Sms user channel", 30030) }
	FailedChooseMailChannel      = func() (*common.BusinessError) { return common.NewBusinessError("Failed to choose mail channel", 30031) }
	FailedChooseSmsExpireChannel = func() (*common.BusinessError) { return common.NewBusinessError("Failed to choose sms expire channel", 30032) }
	FailedMatchPhoneNumber       = func() (*common.BusinessError) { return common.NewBusinessError("Failed to match phone number", 30033) }
	FailedEnqueueExpireMessageId = func() (*common.BusinessError) { return common.NewBusinessError("Failed to enqueue expire message id", 30034) }
	FailedModifySmsMessageId     = func() (*common.BusinessError) { return common.NewBusinessError("Failed to modify Sms message id", 30035) }
	SmsParameterContainComma     = func() (*common.BusinessError) { return common.NewBusinessError("Sms parameter contains [,]", 30036) }
	InvalidSmsParameterCount     = func() (*common.BusinessError) { return common.NewBusinessError("Invalid sms parameter's count", 30037) }
	SystemAliasBlank             = func() (*common.BusinessError) { return common.NewBusinessError("[SystemAlias] is blank", 30038) }
	InvalidTemplateId            = func() (*common.BusinessError) { return common.NewBusinessError("Invalid template id", 30039) }
	NotSendSmsPermission         = func() (*common.BusinessError) { return common.NewBusinessError("No send Sms permission", 30040) }
	NotSendMailPermission        = func() (*common.BusinessError) { return common.NewBusinessError("No send mail permission", 30041) }
	NullModifyParameter          = func() (*common.BusinessError) { return common.NewBusinessError("Null modify parameter", 30042) }
	SystemAliasNotExist          = func() (*common.BusinessError) { return common.NewBusinessError("[SystemAlias] is not exist", 30043) }
	FailedSendSmsMessage         = func() (*common.BusinessError) { return common.NewBusinessError("Failed to send Sms message", 30044) }
)
