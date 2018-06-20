package global

import (
	"fmt"
	"encoding/json"
	"rz/middleware/notifycenter/common"
	"errors"
)

var (
	Config = getConfiguration()
)

type Configuration struct {
	ConsumingInterval int      `json:"consumingInterval"`
	Web               web      `json:"web"`
	Redis             redis    `json:"redis"`
	Sms               sms      `json:"sms"`
	Mail              mail     `json:"mail"`
	QYWeixin          qyWeixin `json:"qyWeixin"`
}

type web struct {
	Listen string `json:"listen"`
}

type redis struct {
	Address    string `json:"address"`
	DatabaseId int    `json:"databaseId"`
	Password   string `json:"password"`
	Master     string `json:"master"`
}

type sms struct {
	Url               string `json:"url"`
	AppKey            string `json:"appKey"`
	AppId             string `json:"appId"`
	DefaultNationCode string `json:"defaultNationCode"`
}

type mail struct {
	Host        string
	Port        int
	UserName    string
	Password    string
	From        string
	ContentType string
}

type qyWeixin struct {
	CorpId     string `json:"corpId"`
	CorpSecret string `json:"corpSecret"`
	AgentId    int    `json:"agentId"`
}

func getConfiguration() (Configuration) {
	filePath := Arguments[ArgumentNameConfig]

	if !common.IsExistPath(filePath) {
		panic(errors.New(fmt.Sprintf("cannot find config file path(%s)\n", filePath)))
	}

	content, err := common.ReadFileContent(filePath)
	if nil != err {
		panic(errors.New(fmt.Sprintf("failed to get config file content(%s). error: %s\n", content, err.Error())))
	}

	var configuration Configuration
	err = json.Unmarshal([]byte(content), &configuration)
	if nil != err {
		panic(errors.New(fmt.Sprintf("invaild config file content(%s). error: %s", content, err.Error())))
	}

	return configuration
}
