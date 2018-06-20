package controllers

import (
	"reflect"

	"rz/middleware/notifycenter/common"
	"rz/middleware/notifycenter/global"
)

type ControllerBase struct {
}

func (myself *ControllerBase) Enable(controller interface{}, isStandard bool) {
	messageControllerType := reflect.ValueOf(controller)
	fieldCount := messageControllerType.NumField()

	for i := 0; i < fieldCount; i++ {
		field := messageControllerType.Field(i)
		// CanInterface check the field is public or private
		if field.CanInterface() {
			controllerPack, ok := field.Interface().(*common.ControllerPack)
			if ok {
				if isStandard {
					global.WebService.RegisterStandardController(controllerPack)
				} else {
					global.WebService.RegisterCommonController(controllerPack)
				}
			}
		}
	}
}
