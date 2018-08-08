package common

import "fmt"

func NewBusinessError(defaultMessage string, code int) (*BusinessError) {
	return &BusinessError{
		DefaultMessage: defaultMessage,
		Code:           code,
	}
}

// implementation of error
type BusinessError struct {
	DefaultMessage string
	message        string
	Code           int
	rawError       error
}

func (myself *BusinessError) AttachError(rawError error) (*BusinessError) {
	myself.rawError = rawError

	return myself
}

func (myself *BusinessError) AttachMessage(message interface{}) (*BusinessError) {
	myself.message = fmt.Sprint(message)

	return myself
}

func (myself *BusinessError) Error() (string) {
	errorMessage := myself.DefaultMessage
	if nil != myself.rawError {
		errorMessage += ". raw error: " + myself.rawError.Error()
	}
	if "" != myself.message {
		errorMessage += ". message: " + myself.message
	}

	return errorMessage
}
