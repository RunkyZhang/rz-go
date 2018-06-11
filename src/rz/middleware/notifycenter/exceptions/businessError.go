package exceptions

import "fmt"

func newBusinessError(message string, code int) (*BusinessError) {
	return &BusinessError{
		Message: message,
		Code:    code,
	}
}

// implementation of error
type BusinessError struct {
	Message string
	Code    int
	Err     error
}

func (businessError *BusinessError) AttachError(err error) (error) {
	businessError.Err = err

	return businessError
}

func (businessError *BusinessError) Error() (string) {
	if nil != businessError.Err {
		return fmt.Sprintf("%s. error: %s", businessError.Message, businessError.Err.Error())
	}

	return businessError.Message
}
