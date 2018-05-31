package global

import (
	"os"
	"fmt"
	"encoding/json"
	"rz/middleware/notifycenter/common"
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
	DefaultNationCode string    `json:"defaultNationCode"`
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
		fmt.Println("Cannot find config file path[", filePath, "].")
		os.Exit(0)
	}

	content, err := common.ReadFileContent(filePath)
	if nil != err {
		fmt.Println("Failed to get config file content[", content, "]. error: ", err.Error(), ".")
		os.Exit(0)
	}

	var configuration Configuration
	err = json.Unmarshal([]byte(content), &configuration)
	if nil != err {
		fmt.Println("Invaild config file content[", content, "]. error: ", err.Error(), ".")
		os.Exit(0)
	}

	return configuration
}
