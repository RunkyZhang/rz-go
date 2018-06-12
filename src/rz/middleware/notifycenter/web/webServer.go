package web

import (
	"encoding/json"
	"net/http"
	s_context "context"
	"time"
	"io"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
	"rz/middleware/notifycenter/global"
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

var (
	server *http.Server
)

func RegisterController(controllerPack *ControllerPack) {
	http.HandleFunc(controllerPack.Pattern, func(responseWriter http.ResponseWriter, request *http.Request) {
		var requestDto models.ResponseDto

		dto, err := controllerPack.ConvertToDtoFunc(request.Body)
		if nil != err {
			requestDto = toResponseDto(err)
			wrapResponseWriter(responseWriter, &requestDto)

			return
		}

		result, err := controllerPack.ControllerFunc(dto)
		if nil != err {
			requestDto = toResponseDto(err)
			wrapResponseWriter(responseWriter, &requestDto)

			return
		}

		requestDto = models.ResponseDto{
			Code:    exceptions.Ok.Code,
			Message: exceptions.Ok.Message,
			Data:    result,
		}
		wrapResponseWriter(responseWriter, &requestDto)
	})
}

func toResponseDto(err error) models.ResponseDto {
	businessError, ok := err.(*exceptions.BusinessError)
	if ok {
		return models.ResponseDto{
			Code:    businessError.Code,
			Message: businessError.Message,
			Data:    nil,
		}
	}

	return models.ResponseDto{
		Code:    exceptions.InternalServerError.Code,
		Message: fmt.Sprintf("%s. error: %s", exceptions.InternalServerError.Message, err.Error()),
		Data:    nil,
	}
}

func Start() {
	go start()
}

func Stop() (error) {
	if nil == server {
		return errors.New("the server is not started")
	}

	context, _ := s_context.WithTimeout(s_context.Background(), 5*time.Second)
	return server.Shutdown(context)
}

func wrapResponseWriter(responseWriter http.ResponseWriter, requestDto *models.ResponseDto) {
	bytes, err := json.Marshal(requestDto)
	if nil != err {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Add("Content-Type", "application/json;charset=UTF-8")
	responseWriter.Write(bytes)
}

func start() (error) {
	http.HandleFunc("/health", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Write([]byte("ok"))
	})

	server = &http.Server{
		Addr: global.Config.Web.Listen,
		// 1 << 10 = 1024, 1 << 20 = 1024 * 1024
		MaxHeaderBytes: 1 << 20,
	}
	server.SetKeepAlivesEnabled(true)
	return server.ListenAndServe()
}
