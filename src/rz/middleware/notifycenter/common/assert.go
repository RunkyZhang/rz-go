package common

import "errors"

var (
	Assert = &assert{}
)

type assert struct {
}

func (*assert) IsNotNil(value interface{}, name string) {
	if nil == value {
		panic("the parameter[" + name + "] is nil")
	}
}

func (*assert) IsNilError(err error, message string) {
	if nil != err {
		if "" == message {
			panic("the err[" + err.Error() + "] is nil")
		} else {
			panic("the err[" + err.Error() + "] is nil. message: " + message)
		}
	}
}

func (*assert) NewNilParameterError(name string) (error) {
	return errors.New("the parameter[" + name + "] is nil")
}
