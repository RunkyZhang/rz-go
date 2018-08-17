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

func (*assert) IsNilErrorToPanic(err error, message string) {
	if nil != err {
		if "" == message {
			panic(errors.New(fmt.Sprintf("The err(%s) is nil", err.Error())))
		} else {
			panic(errors.New(fmt.Sprintf("The err(%s) is nil. message: %s", err.Error(), message)))
		}
	}
}

func (*assert) IsNotBlankToError(value string, name string) (error) {
	if "" == value {
		return errors.New(fmt.Sprintf("The parameter(%s) is blank", name))
	}

	return nil
}

func (*assert) IsNotBlankToPanic(value string, name string) {
	if "" == value {
		panic(errors.New(fmt.Sprintf("The parameter(%s) is blank", name)))
	}
}

func (*assert) IsTrueToError(ok bool, expression string) (error) {
	if !ok {
		return errors.New(fmt.Sprintf("The expression(%s) is not true", expression))
	}

	return nil
}

func (*assert) IsTrueToPanic(ok bool, expression string) {
	if !ok {
		panic(errors.New(fmt.Sprintf("The expression(%s) is not true", expression)))
	}
}
