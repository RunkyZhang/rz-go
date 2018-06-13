package exceptions

type newBusinessErrorFunc func() (*BusinessError)

var (
	Ok                      = func() (*BusinessError) { return newBusinessError("Ok", 0) }
	InternalServerError     = func() (*BusinessError) { return newBusinessError("Internal server error", 30000) }
	DtoNull                 = func() (*BusinessError) { return newBusinessError("[Dto] is null", 30001) }
	InvalidDtoType          = func() (*BusinessError) { return newBusinessError("Invalid [Dto] type", 30002) }
	ErrorTosEmpty           = func() (*BusinessError) { return newBusinessError("[Tos] is empty", 30003) }
	SubjectBlank            = func() (*BusinessError) { return newBusinessError("The subject is blank", 30004) }
	InvalidSendChannel      = func() (*BusinessError) { return newBusinessError("Invalid send channel number", 30005) }
	InvalidMessageState     = func() (*BusinessError) { return newBusinessError("Invalid message state", 30006) }
	MessageExpire           = func() (*BusinessError) { return newBusinessError("Message expire", 30007) }
	MessageBodyMissed       = func() (*BusinessError) { return newBusinessError("Message body is missed", 30008) }
	FailedAddSmsUserMessage = func() (*BusinessError) { return newBusinessError("Failed to add sms user message", 30009) }
	FailedEnqueueMessageId  = func() (*BusinessError) { return newBusinessError("Failed to enqueue message id", 30010) }
	InvalidExtend           = func() (*BusinessError) { return newBusinessError("Invalid extend", 30011) }
	TemplateIdNotExist      = func() (*BusinessError) { return newBusinessError("Template id is not exist", 30012) }
	InvalidPattern          = func() (*BusinessError) { return newBusinessError("Invalid pattern", 30013) }
	PatternNotMatch         = func() (*BusinessError) { return newBusinessError("Pattern is not match string", 30014) }
)
