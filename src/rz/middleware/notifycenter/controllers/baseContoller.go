package controllers

import (
	"reflect"

	"rz/middleware/notifycenter/web"
)

type BaseController struct {
}

func (*BaseController) Enable(controller interface{}) {
	messageControllerType := reflect.ValueOf(controller)
	fieldCount := messageControllerType.NumField()

	for i := 0; i < fieldCount; i++ {
		field := messageControllerType.Field(i)
		// CanInterface check the field is public or private
		if field.CanInterface() {
			controllerPack, ok := field.Interface().(*web.ControllerPack)
			if ok {
				web.RegisterController(controllerPack)
			}
		}
	}
}
