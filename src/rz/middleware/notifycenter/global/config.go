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
	Web      web      `json:"web"`
	Redis    redis    `json:"redis"`
	QYWeixin qyWeixin `json:"qyWeixin"`
}

type web struct {
	Listen string `json:"listen"`
}

type redis struct {
	Addresses  []string `json:"addresses"`
	DatabaseId int      `json:"databaseId"`
	Password   string   `json:"password"`
	Master     string   `json:"master"`
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
