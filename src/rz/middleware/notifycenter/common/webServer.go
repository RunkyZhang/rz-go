package common

import (
	"encoding/json"
	"net/http"
	"context"
	"time"
	"io"

	"rz/middleware/notifycenter/exceptions"
	"errors"
	"fmt"
)

type ConvertToDtoFunc func(body io.ReadCloser) (interface{}, error)

type ControllerFunc func(interface{}) (interface{}, error)

type ControllerPack struct {
	Pattern          string
	ControllerFunc   ControllerFunc
	ConvertToDtoFunc ConvertToDtoFunc
}

type ResponseDto struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type HealthReport struct {
	Ok      bool                   `json:"ok"`
	Name    string                 `json:"name"`
	Message string                 `json:"message"`
	Type    string                 `json:"type"`
	Level   int                    `json:"level"`
	Detail  map[string]interface{} `json:"detail"`
}

type HealthIndicator interface {
	Indicate() (*HealthReport)
}

func NewWebService(address string) (*webService) {
	webService := &webService{
		address:          address,
		healthIndicators: []HealthIndicator{},
	}

	return webService
}

type webService struct {
	server           *http.Server
	address          string
	healthIndicators []HealthIndicator
}

func (myself *webService) RegisterStandardController(controllerPack *ControllerPack) {
	var responseDto ResponseDto

	http.HandleFunc(controllerPack.Pattern, func(responseWriter http.ResponseWriter, request *http.Request) {
		//defer func() {
		//	value := recover()
		//	if nil != value {
		//		responseDto = myself.errorToResponseDto(value)
		//		myself.wrapResponseWriter(responseWriter, &responseDto)
		//	}
		//}()

		dto, err := controllerPack.ConvertToDtoFunc(request.Body)
		if nil != err {
			responseDto = myself.errorToResponseDto(err)
			myself.wrapResponseWriter(responseWriter, &responseDto)

			return
		}

		result, err := controllerPack.ControllerFunc(dto)
		if nil != err {
			responseDto = myself.errorToResponseDto(err)
			myself.wrapResponseWriter(responseWriter, &responseDto)

			return
		}

		exceptionsOk := exceptions.Ok()
		responseDto = ResponseDto{
			Code:    exceptionsOk.Code,
			Message: exceptionsOk.Error(),
			Data:    result,
		}
		myself.wrapResponseWriter(responseWriter, &responseDto)
	})
}

func (myself *webService) RegisterCommonController(controllerPack *ControllerPack) {
	http.HandleFunc(controllerPack.Pattern, func(responseWriter http.ResponseWriter, request *http.Request) {
		//defer func() {
		//	value := recover()
		//	if nil != value {
		//		http.Error(responseWriter, fmt.Sprintln(value), http.StatusInternalServerError)
		//	}
		//}()

		dto, err := controllerPack.ConvertToDtoFunc(request.Body)
		if nil != err {
			http.Error(responseWriter, exceptions.InvalidDtoType().AttachError(err).Error(), http.StatusInternalServerError)

			return
		}

		result, err := controllerPack.ControllerFunc(dto)
		if nil != err {
			http.Error(responseWriter, exceptions.FailedInvokeController().AttachError(err).Error(), http.StatusInternalServerError)

			return
		}

		myself.wrapResponseWriter(responseWriter, result)
	})
}

func (myself *webService) RegisterHealthIndicator(healthIndicator HealthIndicator) {
	if nil == healthIndicator {
		return
	}

	myself.healthIndicators = append(myself.healthIndicators, healthIndicator)
}

func (myself *webService) Start() {
	go myself.start()
}

func (myself *webService) Stop() (error) {
	if nil == myself.server {
		return errors.New("the server is not started")
	}

	timeoutContext, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return myself.server.Shutdown(timeoutContext)
}

func (*webService) errorToResponseDto(value interface{}) ResponseDto {
	businessError, ok := value.(*exceptions.BusinessError)
	if ok {
		return ResponseDto{
			Code:    businessError.Code,
			Message: businessError.Error(),
			Data:    nil,
		}
	}

	exceptionsInternalServerError := exceptions.InternalServerError()
	err, ok := value.(error)
	if ok {
		return ResponseDto{
			Code:    exceptionsInternalServerError.Code,
			Message: exceptionsInternalServerError.AttachError(err).Error(),
			Data:    nil,
		}
	}

	return ResponseDto{
		Code:    exceptionsInternalServerError.Code,
		Message: exceptionsInternalServerError.AttachMessage(fmt.Sprintln(value)).Error(),
		Data:    nil,
	}
}

func (myself *webService) wrapResponseWriter(responseWriter http.ResponseWriter, body interface{}) {
	bytes, err := json.Marshal(body)
	if nil != err {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Add("Content-Type", "application/json;charset=UTF-8")
	responseWriter.Write(bytes)
}

func (myself *webService) health() {
	http.HandleFunc("/health", func(responseWriter http.ResponseWriter, request *http.Request) {
		var healthReports []*HealthReport

		//defer func() {
		//	value := recover()
		//	if nil != value {
		//		healthReport := &HealthReport{
		//			Ok:      false,
		//			Name:    "unknown error",
		//			Message: fmt.Sprintln(value),
		//			Type:    "panic",
		//			Level:   0,
		//		}
		//		healthReports = append(healthReports, healthReport)
		//
		//		myself.wrapResponseWriter(responseWriter, healthReports)
		//	}
		//}()

		length := len(myself.healthIndicators)
		for i := 0; i < length; i++ {
			healthIndicator := myself.healthIndicators[i]
			healthReports = append(healthReports, healthIndicator.Indicate())
		}

		myself.wrapResponseWriter(responseWriter, healthReports)
	})
}

func (myself *webService) start() (error) {
	myself.health()

	myself.server = &http.Server{
		Addr: myself.address,
		// 1 << 10 = 1024, 1 << 20 = 1024 * 1024
		MaxHeaderBytes: 1 << 20,
	}
	myself.server.SetKeepAlivesEnabled(true)
	return myself.server.ListenAndServe()
}
