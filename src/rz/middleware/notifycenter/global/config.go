package global

import (
	"fmt"
	"encoding/json"
	"errors"
	"sync"
	"os"

	"rz/middleware/notifycenter/common"
)

var (
	configuration *Configuration = nil
	configLock    sync.Mutex
)

type Configuration struct {
	ConsumingInterval int               `json:"consumingInterval"`
	Web               web               `json:"web"`
	Redis             redis             `json:"redis"`
	Sms               sms               `json:"sms"`
	Mail              mail              `json:"mail"`
	QYWeixin          qyWeixin          `json:"qyWeixin"`
	Databases         map[string]string `json:"databases"`
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

func GetConfig() (*Configuration) {
	if nil != configuration {
		return configuration
	}

	configLock.Lock()
	defer configLock.Unlock()

	if nil != configuration {
		return configuration
	}

	filePath := Arguments[ArgumentNameConfig]
	environmentId := os.Getenv("CONFIGENV")
	filePath = fmt.Sprintf(filePath, getConfigFileSuffix(environmentId))

	if !common.IsExistPath(filePath) {
		panic(errors.New(fmt.Sprintf("cannot find config file path(%s)\n", filePath)))
	}

	content, err := common.ReadFileContent(filePath)
	if nil != err {
		panic(errors.New(fmt.Sprintf("failed to get config file content(%s). error: %s\n", content, err.Error())))
	}

	configuration = &Configuration{}
	err = json.Unmarshal([]byte(content), &configuration)
	if nil != err {
		panic(errors.New(fmt.Sprintf("invaild config file content(%s). error: %s", content, err.Error())))
	}

	return configuration
}

func getConfigFileSuffix(environmentId string) (string) {
	if "10" == environmentId {
		return "_DEV"
	} else if "20" == environmentId {
		return "_TEST"
	} else if "30" == environmentId {
		return "_UAT"
	} else if "40" == environmentId {
		return "_PRD"
	} else if "50" == environmentId {
		return "_MIT"
	} else if "60" == environmentId {
		return "_STG"
	}

	return "_DEV"
}

func RefreshConfig() {
	oldConfiguration := configuration

	defer func() {
		value := recover()
		if nil != value {
			fmt.Printf("failed to refresh config file\n")

			configuration = oldConfiguration
		}
	}()

	configuration = nil
	GetConfig()
}
