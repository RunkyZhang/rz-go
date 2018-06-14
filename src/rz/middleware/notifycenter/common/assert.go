package common

import (
	"errors"
	"fmt"
)

var (
	Assert = &assert{}
)

type assert struct {
}

func (*assert) IsNotNilToPanic(value interface{}, name string) {
	if nil == value {
		panic("the parameter[" + name + "] is nil")
	}
}

func (*assert) IsNilErrorToPanic(err error, message string) {
	if nil != err {
		if "" == message {
			panic("the err[" + err.Error() + "] is nil")
		} else {
			panic("the err[" + err.Error() + "] is nil. message: " + message)
		}
	}
}

func (*assert) IsNotNilToError(value interface{}, name string) (error) {
	if nil == value {
		return errors.New(fmt.Sprintf("the parameter(%s) is nil", name))
	}

	return nil
}
