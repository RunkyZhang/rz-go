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
		panic(fmt.Sprintf("the parameter(%s) is nil", name))
	}
}

func (*assert) IsNilErrorToPanic(err error, message string) {
	if nil != err {
		if "" == message {
			panic(fmt.Sprintf("the err(%s) is nil", err.Error()))
		} else {
			panic(fmt.Sprintf("the err(%s) is nil. message: %s", err.Error(), message))
		}
	}
}

func (*assert) IsNotNilToError(value interface{}, name string) (error) {
	if nil == value {
		return errors.New(fmt.Sprintf("the parameter(%s) is nil", name))
	}

	return nil
}
