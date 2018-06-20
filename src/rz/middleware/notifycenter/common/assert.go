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
		panic(errors.New(fmt.Sprintf("the parameter(%s) is nil", name)))
	}
}

func (*assert) IsNilErrorToPanic(err error, message string) {
	if nil != err {
		if "" == message {
			panic(errors.New(fmt.Sprintf("the err(%s) is nil", err.Error())))
		} else {
			panic(errors.New(fmt.Sprintf("the err(%s) is nil. message: %s", err.Error(), message)))
		}
	}
}

func (*assert) IsNotNilToError(value interface{}, name string) (error) {
	if nil == value {
		return errors.New(fmt.Sprintf("the parameter(%s) is nil", name))
	}

	return nil
}

func (*assert) IsTrueToError(ok bool, expression string) (error) {
	if !ok {
		return errors.New(fmt.Sprintf("the expression(%s) is not true", expression))
	}

	return nil
}
