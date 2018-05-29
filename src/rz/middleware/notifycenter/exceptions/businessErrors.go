package exceptions

var (
	Ok                  = newBusinessError("Ok", 0)
	InternalServerError = newBusinessError("Internal server error", 30000)
	DtoNull             = newBusinessError("[Dto] is null", 30001)
	InvalidDtoType      = newBusinessError("Invalid [Dto] type", 30002)
	ErrorTosEmpty       = newBusinessError("[Tos] is empty", 30003)
	SubjectBlank        = newBusinessError("The subject is blank", 30004)
	InvalidSendChannel  = newBusinessError("Invalid send channel number", 30005)
	InvalidMessageState = newBusinessError("Invalid message state", 30006)
	MessageExpire       = newBusinessError("Message expire", 30007)
	MessageBodyMissed       = newBusinessError("Message body is missed", 30008)
)
