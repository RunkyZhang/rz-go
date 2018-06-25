package common

import (
	"encoding/json"
	"net/http"
	"context"
	"time"
	"errors"
	"fmt"
	"bytes"
	"math/rand"
)

type ConvertToDtoFunc func(body []byte) (interface{}, error)

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
	http.HandleFunc(controllerPack.Pattern, func(responseWriter http.ResponseWriter, request *http.Request) {
		id := myself.buildRequestId()

		defer func() {
			value := recover()
			if nil != value {
				responseDto := myself.errorToResponseDto(value)
				myself.wrapResponseWriter(responseWriter, request, id, &responseDto, nil, "")
			}
		}()

		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(request.Body)
		if nil != err {
			responseDto := myself.errorToResponseDto(err)
			myself.wrapResponseWriter(responseWriter, request, id, &responseDto, nil, "")

			return
		}

		myself.log("Start", id, request.URL.String(), request.Method, buffer.Bytes())

		dto, err := controllerPack.ConvertToDtoFunc(buffer.Bytes())
		if nil != err {
			responseDto := myself.errorToResponseDto(err)
			myself.wrapResponseWriter(responseWriter, request, id, &responseDto, nil, "")

			return
		}

		result, err := controllerPack.ControllerFunc(dto)
		if nil != err {
			responseDto := myself.errorToResponseDto(err)
			myself.wrapResponseWriter(responseWriter, request, id, &responseDto, nil, "")

			return
		}

		responseDto := ResponseDto{
			Code:    0,
			Message: "Ok",
			Data:    result,
		}
		myself.wrapResponseWriter(responseWriter, request, id, &responseDto, nil, "")
	})
}

func (myself *webService) RegisterCommonController(controllerPack *ControllerPack) {
	http.HandleFunc(controllerPack.Pattern, func(responseWriter http.ResponseWriter, request *http.Request) {
		id := myself.buildRequestId()

		defer func() {
			value := recover()
			if nil != value {
				myself.wrapResponseWriter(responseWriter, request, id, nil, value, "failed by panic")
			}
		}()

		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(request.Body)
		if nil != err {
			myself.wrapResponseWriter(responseWriter, request, id, nil, err.Error(), "failed to read bytes from body")
			return
		}

		myself.log("Start", id, request.URL.String(), request.Method, buffer.Bytes())

		dto, err := controllerPack.ConvertToDtoFunc(buffer.Bytes())
		if nil != err {
			myself.wrapResponseWriter(responseWriter, request, id, nil, err.Error(), "failed to convert body to [Dto]")
			return
		}

		result, err := controllerPack.ControllerFunc(dto)
		if nil != err {
			myself.wrapResponseWriter(responseWriter, request, id, nil, err.Error(), "failed to invoke controller")
			return
		}

		myself.wrapResponseWriter(responseWriter, request, id, result, nil, "")
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

func (myself *webService) health() {
	http.HandleFunc("/health", func(responseWriter http.ResponseWriter, request *http.Request) {
		id := myself.buildRequestId()
		var healthReports []*HealthReport

		defer func() {
			value := recover()
			if nil != value {
				healthReport := &HealthReport{
					Ok:      false,
					Name:    "unknown error",
					Message: fmt.Sprintln(value),
					Type:    "panic",
					Level:   0,
				}
				healthReports = append(healthReports, healthReport)

				myself.wrapResponseWriter(responseWriter, request, id, healthReports, nil, "")
			}
		}()

		myself.log("Start", id, request.URL.String(), request.Method, nil)

		length := len(myself.healthIndicators)
		for i := 0; i < length; i++ {
			healthIndicator := myself.healthIndicators[i]
			healthReports = append(healthReports, healthIndicator.Indicate())
		}

		myself.wrapResponseWriter(responseWriter, request, id, healthReports, nil, "")
	})
}

func (*webService) errorToResponseDto(value interface{}) ResponseDto {
	businessError, ok := value.(*BusinessError)
	if ok {
		return ResponseDto{
			Code:    businessError.Code,
			Message: businessError.Error(),
			Data:    nil,
		}
	}

	err, ok := value.(error)
	if ok {
		return ResponseDto{
			Code:    1,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return ResponseDto{
		Code:    1,
		Message: fmt.Sprintln(value),
		Data:    nil,
	}
}

func (myself *webService) wrapResponseWriter(responseWriter http.ResponseWriter, request *http.Request, id string, body interface{}, errorValue interface{}, message string) {
	errorMessage := ""
	if nil != errorValue {
		err, ok := errorValue.(error)
		if ok {
			errorMessage = fmt.Sprintf("%s; error: %s", message, err.Error())
			return
		} else {
			errorMessage = fmt.Sprintf("%s; error: %s", message, errorValue)
		}
	}
	var buffer []byte
	buffer, err := json.Marshal(body)
	if nil != err {
		errorMessage = fmt.Sprintf("failed to convert body to json; error: %s", err.Error())
	}

	if "" != errorMessage {
		myself.log("Failed", id, request.URL.String(), request.Method, []byte(errorMessage))
		http.Error(responseWriter, errorMessage, http.StatusInternalServerError)
		return
	}

	responseDto, ok := body.(*ResponseDto)
	if ok && 0 != responseDto.Code {
		myself.log("Failed", id, request.URL.String(), request.Method, buffer)
	} else {
		myself.log("Success", id, request.URL.String(), request.Method, buffer)
	}

	responseWriter.Header().Add("Content-Type", "application/json;charset=UTF-8")
	responseWriter.Write(buffer)
}

func (myself *webService) log(title string, id string, url string, method string, body []byte) {
	fmt.Printf("%s-[%s][%s][%s][body: %s]\n", title, id, url, method, body)
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

func (myself *webService) buildRequestId() (string) {
	var randomNumber = Int32ToString(rand.Intn(10000))

	return Int64ToString(time.Now().Unix()) + "-" + randomNumber
}
