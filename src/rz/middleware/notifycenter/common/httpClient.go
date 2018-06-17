package common

//import "git.zhaogangren.com/cloud/cloud.base.utils-go.sdk/httplib"

func NewHttpClient() (*HttpClient) {
	httpClient := &HttpClient{}

	return httpClient
}

type HttpClient struct {
}

func (myself *HttpClient) Get(uri string) ([]byte, error) {
	return nil, nil
}

func (myself *HttpClient) Post(uri string, body interface{}) ([]byte, error) {
	//return httplib.Post(uri, body)

	return nil, nil
}

func (myself *HttpClient) Put(uri string, body interface{}) ([]byte, error) {
	return nil, nil
}

func (myself *HttpClient) Delete(uri string) ([]byte, error) {
	return nil, nil
}
