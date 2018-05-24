package exceptions

import "rz/middleware/notifycenter/models"

var (
	Ok                  = newBusinessError("Ok", 0)
	InternalServerError = newBusinessError("Internal server error", 30000)
	DtoNull             = newBusinessError("[Dto] is null", 30001)
	InvalidDtoType      = newBusinessError("Invalid [Dto] type", 30002)
	ErrorTosEmpty       = newBusinessError("[Tos] is empty", 30003)
	SubjectBlank        = newBusinessError("The subject is blank", 30004)
	InvalidSendChannel  = newBusinessError("Invalid send channel number", 30005)
)

func ToResponseDto(err error) models.ResponseDto {
	businessError, ok := err.(*BusinessError)
	if ok {
		return models.ResponseDto{
			Code:    businessError.Code,
			Message: businessError.Message,
			Data:    nil,
		}
	} else {
		return models.ResponseDto{
			Code:    InternalServerError.Code,
			Message: InternalServerError.Message,
			Data:    nil,
		}
	}
}
