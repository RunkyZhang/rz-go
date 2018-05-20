package web

import (
	"encoding/json"
	"net/http"
	s_context "context"
	"time"
	"fmt"
	"io"

	"rz/middleware/notifycenter/models"
	"rz/middleware/notifycenter/exceptions"
)

type ConvertToDtoFunc func(body io.ReadCloser) (interface{}, error)

type ControllerFunc func(interface{}) (interface{}, error)

type VerifyFunc func(interface{}) (error)

type ControllerPack struct {
	Pattern          string
	ControllerFunc   ControllerFunc
	ConvertToDtoFunc ConvertToDtoFunc
	VerifyFunc       VerifyFunc
}

var (
	server *http.Server
)

func RegisterController(controllerPack *ControllerPack) {
	http.HandleFunc(controllerPack.Pattern, func(responseWriter http.ResponseWriter, request *http.Request) {
		var requestDto models.ResponseDto

		dto, exception := controllerPack.ConvertToDtoFunc(request.Body)
		if nil != exception {
			requestDto = exceptions.ToResponseDto(exception)
			wrapResponseWriter(responseWriter, &requestDto)

			return
		}

		exception = controllerPack.VerifyFunc(dto)
		if nil != exception {
			requestDto = exceptions.ToResponseDto(exception)
			wrapResponseWriter(responseWriter, &requestDto)

			return
		}

		result, exception := controllerPack.ControllerFunc(dto)
		if nil != exception {
			requestDto = exceptions.ToResponseDto(exception)
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

func Start() {
	go start()
}

func Stop() {
	if nil == server {
		return
	}

	context, _ := s_context.WithTimeout(s_context.Background(), 5*time.Second)
	exception := server.Shutdown(context)

	fmt.Println("failed to shutdown web server: ", exception)
}

func wrapResponseWriter(responseWriter http.ResponseWriter, requestDto *models.ResponseDto) {
	bytes, exception := json.Marshal(requestDto)
	if nil != exception {
		http.Error(responseWriter, exception.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Add("Content-Type", "application/json;charset=UTF-8")
	responseWriter.Write(bytes)
}

func start() {
	http.HandleFunc("/health", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Write([]byte("ok"))
	})

	server = &http.Server{
		Addr: "0.0.0.0:3030",
		// 1 << 10 = 1024, 1 << 20 = 1024 * 1024
		MaxHeaderBytes: 1 << 20,
	}
	server.SetKeepAlivesEnabled(true)
	fmt.Println(server.ListenAndServe())
}
