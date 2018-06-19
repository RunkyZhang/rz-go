package controllers

import (
	"reflect"

	"rz/middleware/notifycenter/global"
	"rz/middleware/notifycenter/common"
)

type ControllerBase struct {
}

func (myself *ControllerBase) Enable(controller interface{}) {
	messageControllerType := reflect.ValueOf(controller)
	fieldCount := messageControllerType.NumField()

	for i := 0; i < fieldCount; i++ {
		field := messageControllerType.Field(i)
		// CanInterface check the field is public or private
		if field.CanInterface() {
			controllerPack, ok := field.Interface().(*common.ControllerPack)
			if ok {
				global.WebService.RegisterStandardController(controllerPack)
			}
		}
	}
}
