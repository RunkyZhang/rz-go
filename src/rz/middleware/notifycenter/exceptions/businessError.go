package exceptions

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
}

func (businessError *BusinessError) Error() (string) {
	return businessError.Message
}
