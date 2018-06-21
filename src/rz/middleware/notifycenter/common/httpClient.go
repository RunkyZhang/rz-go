package common

import (
	"net/http"
	"encoding/json"
	"bytes"
	"strings"
	"compress/gzip"
	"io/ioutil"
	"time"
	"context"
	"net"
	"net/url"
)

var DefaultHttpClientSettings = HttpClientSettings{
	ConnectTimeout: 3 * time.Second,
	RequestTimeout: 100 * time.Second,
	KeepAliveTime:  30 * time.Second,
	AcceptContent:  "application/json;charset=utf-8",
	AcceptEncoding: "gzip",
}

func NewHttpClient(httpClientSettings *HttpClientSettings) (*HttpClient) {
	httpClient := &HttpClient{
		headers: map[string]string{},
	}
	if nil == httpClientSettings {
		httpClient.httpClientSettings = DefaultHttpClientSettings
	} else {
		httpClient.httpClientSettings = *httpClientSettings
	}

	httpClient.headers["Accept-Encoding"] = httpClient.httpClientSettings.AcceptEncoding
	httpClient.headers["Accept"] = httpClient.httpClientSettings.AcceptContent
	httpClient.headers["Connection"] = "Keep-Alive"
	httpClient.headers["Content-Type"] = "application/json;charset=utf-8"

	return httpClient
}

type HttpClientSettings struct {
	ConnectTimeout time.Duration
	RequestTimeout time.Duration
	KeepAliveTime  time.Duration
	AcceptContent  string
	AcceptEncoding string
}

type HttpClient struct {
	baseUrl            string
	headers            map[string]string
	httpClientSettings HttpClientSettings
}

func (myself *HttpClient) SetBaseUrl(baseUrl string) (*HttpClient) {
	if "" != baseUrl && strings.HasSuffix(baseUrl, "/") {
		myself.baseUrl = baseUrl[0 : len(baseUrl)-1]
	}

	return myself
}

func (myself *HttpClient) SetAcceptContent(acceptContent string) (*HttpClient) {
	myself.httpClientSettings.AcceptContent = acceptContent
	myself.headers["Accept"] = myself.httpClientSettings.AcceptContent

	return myself
}

func (myself *HttpClient) SetAcceptEncoding(acceptEncoding string) (*HttpClient) {
	myself.httpClientSettings.AcceptEncoding = acceptEncoding
	myself.headers["Accept-Encoding"] = myself.httpClientSettings.AcceptEncoding

	return myself
}

func (myself *HttpClient) SetHeaders(headers map[string]string) (*HttpClient) {
	if nil == headers {
		return myself
	}

	for key, value := range headers {
		myself.headers[key] = value
	}

	return myself
}

func (myself *HttpClient) SetTimeout(connectTimeout time.Duration, requestTimeout time.Duration, keepAliveTime time.Duration) (*HttpClient) {
	myself.httpClientSettings.ConnectTimeout = connectTimeout
	myself.httpClientSettings.RequestTimeout = requestTimeout
	myself.httpClientSettings.KeepAliveTime = keepAliveTime

	return myself
}

func (myself *HttpClient) Get(uri string, headers ...map[string]string) ([]byte, error) {
	return myself.request("GET", uri, nil, headers...)
}

func (myself *HttpClient) Post(uri string, body interface{}, headers ...map[string]string) ([]byte, error) {
	//return httplib.Post(uri, body)

	return myself.request("POST", uri, body, headers...)
}

func (myself *HttpClient) Put(uri string, body interface{}, headers ...map[string]string) ([]byte, error) {
	return myself.request("PUT", uri, body, headers...)
}

func (myself *HttpClient) Delete(uri string, headers ...map[string]string) ([]byte, error) {
	return myself.request("DELETE", uri, nil, headers...)
}

func (myself *HttpClient) request(method string, uri string, body interface{}, headers ...map[string]string) ([]byte, error) {
	var err error
	request := &http.Request{
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	request.URL, err = myself.buildUrl(uri)
	if nil != err {
		return nil, err
	}
	myself.setHeaders(request, headers...)
	myself.setBody(request, body)

	client := myself.buildClient()
	response, err := client.Do(request)
	if nil != err {
		return nil, err
	}
	if nil != response.Body {
		defer response.Body.Close()
	}

	return myself.getBytes(response)
}

func (myself *HttpClient) buildClient() (*http.Client) {
	transport := &http.Transport{
		IdleConnTimeout:       45 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		DisableKeepAlives:     false,
		DisableCompression:    false,
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConnsPerHost:   200,
		MaxIdleConns:          2000,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   myself.httpClientSettings.ConnectTimeout,
				KeepAlive: myself.httpClientSettings.KeepAliveTime,
				DualStack: true,
			}
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   myself.httpClientSettings.RequestTimeout,
	}
}

func (myself *HttpClient) buildUrl(uri string) (*url.URL, error) {
	if "" == myself.baseUrl {
		return url.Parse(uri)
	}

	if "" == uri || strings.HasSuffix(uri, "/") {
		return url.Parse(myself.baseUrl + uri)
	}

	return url.Parse(myself.baseUrl + "/" + uri)
}

func (myself *HttpClient) setBody(request *http.Request, body interface{}) (error) {
	var err error
	var buffer []byte
	if nil != body {
		buffer, err = json.Marshal(body)
		if nil != err {
			return err
		}
	}
	request.Body = ioutil.NopCloser(bytes.NewReader(buffer))
	request.ContentLength = int64(len(buffer))

	return nil
}

func (myself *HttpClient) setHeaders(request *http.Request, headers ...map[string]string) {
	for key, value := range myself.headers {
		request.Header.Set(key, value)
	}

	if nil == headers || 0 == len(headers) {
		return
	}

	keyValues := headers[0]
	for key, value := range keyValues {
		request.Header.Set(key, value)
	}
}

func (myself *HttpClient) getBytes(response *http.Response) ([]byte, error) {
	if strings.EqualFold("gzip", response.Header.Get("Content-Encoding")) {
		reader, err := gzip.NewReader(response.Body)
		if nil != err {
			return nil, err
		}

		return ioutil.ReadAll(reader)
	}

	return ioutil.ReadAll(response.Body)
}
