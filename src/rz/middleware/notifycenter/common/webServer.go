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
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewWebService(address string) (*webService) {
	webService := &webService{
		address: address,
	}

	return webService
}

type webService struct {
	server  *http.Server
	address string
}

func (myself *webService) RegisterStandardController(controllerPack *ControllerPack) {
	var responseDto ResponseDto

	http.HandleFunc(controllerPack.Pattern, func(responseWriter http.ResponseWriter, request *http.Request) {
		defer func() {
			err := recover().(error)
			if nil == err {
				http.Error(responseWriter, "unknown error", http.StatusInternalServerError)
				return
			}

			responseDto = myself.toResponseDto(err)
			myself.wrapResponseWriter(responseWriter, &responseDto)
		}()

		dto, err := controllerPack.ConvertToDtoFunc(request.Body)
		if nil != err {
			responseDto = myself.toResponseDto(err)
			myself.wrapResponseWriter(responseWriter, &responseDto)

			return
		}

		result, err := controllerPack.ControllerFunc(dto)
		if nil != err {
			responseDto = myself.toResponseDto(err)
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
		defer func() {
			err := recover().(error)
			errorMessage := fmt.Sprintf("unknown error: %s", err)

			http.Error(responseWriter, errorMessage, http.StatusInternalServerError)
		}()

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

func (*webService) toResponseDto(err error) ResponseDto {
	businessError, ok := err.(*exceptions.BusinessError)
	if ok {
		return ResponseDto{
			Code:    businessError.Code,
			Message: businessError.Error(),
			Data:    nil,
		}
	}

	exceptionsInternalServerError := exceptions.InternalServerError().AttachError(err)
	return ResponseDto{
		Code:    exceptionsInternalServerError.Code,
		Message: exceptionsInternalServerError.Error(),
		Data:    nil,
	}
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

func (myself *webService) wrapResponseWriter(responseWriter http.ResponseWriter, body interface{}) {
	bytes, err := json.Marshal(body)
	if nil != err {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Add("Content-Type", "application/json;charset=UTF-8")
	responseWriter.Write(bytes)
}

func (myself *webService) start() (error) {
	http.HandleFunc("/health", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Write([]byte("ok"))
	})

	myself.server = &http.Server{
		Addr: myself.address,
		// 1 << 10 = 1024, 1 << 20 = 1024 * 1024
		MaxHeaderBytes: 1 << 20,
	}
	myself.server.SetKeepAlivesEnabled(true)
	return myself.server.ListenAndServe()
}
