package exceptions

import "rz/middleware/notifycenter/common"

type newBusinessErrorFunc func() (*common.BusinessError)

var (
	Ok                       = func() (*common.BusinessError) { return common.NewBusinessError("Ok", 0) }
	InternalServerError      = func() (*common.BusinessError) { return common.NewBusinessError("Internal server error", 30000) }
	DtoNull                  = func() (*common.BusinessError) { return common.NewBusinessError("[Dto] is null", 30001) }
	InvalidDtoType           = func() (*common.BusinessError) { return common.NewBusinessError("Invalid [Dto] type", 30002) }
	TosEmpty                 = func() (*common.BusinessError) { return common.NewBusinessError("[Tos] is empty", 30003) }
	SubjectBlank             = func() (*common.BusinessError) { return common.NewBusinessError("The subject is blank", 30004) }
	InvalidSendChannel       = func() (*common.BusinessError) { return common.NewBusinessError("Invalid send channel number", 30005) }
	InvalidMessageState      = func() (*common.BusinessError) { return common.NewBusinessError("Invalid message state", 30006) }
	MessageExpire            = func() (*common.BusinessError) { return common.NewBusinessError("Message expire", 30007) }
	MessageBodyMissed        = func() (*common.BusinessError) { return common.NewBusinessError("Message body is missed", 30008) }
	FailedAddSmsUserMessage  = func() (*common.BusinessError) { return common.NewBusinessError("Failed to add sms user message", 30009) }
	FailedEnqueueMessageId   = func() (*common.BusinessError) { return common.NewBusinessError("Failed to enqueue message id", 30010) }
	InvalidExtend            = func() (*common.BusinessError) { return common.NewBusinessError("Invalid extend", 30011) }
	TemplateIdNotExist       = func() (*common.BusinessError) { return common.NewBusinessError("Template id is not exist", 30012) }
	InvalidPattern           = func() (*common.BusinessError) { return common.NewBusinessError("Invalid pattern", 30013) }
	PatternNotMatch          = func() (*common.BusinessError) { return common.NewBusinessError("Pattern is not match string", 30014) }
	FailedInvokeController   = func() (*common.BusinessError) { return common.NewBusinessError("Failed invoke controller", 30015) }
	InvalidSystemAlias       = func() (*common.BusinessError) { return common.NewBusinessError("Invalid [SystemAlias]", 30016) }
	InvalidMessageExpireTime = func() (*common.BusinessError) { return common.NewBusinessError("Invalid message expire time", 30017) }
	InvalidIdentifyingCode   = func() (*common.BusinessError) { return common.NewBusinessError("Invalid identifying code", 30018) }
)
