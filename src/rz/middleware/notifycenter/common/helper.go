package common

import (
	"strings"
	s_bytes "bytes"
	"io"
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

func IsStringBlank(value string) (bool) {
	return 0 == len(strings.TrimSpace(value))
}

func ReaderToString(reader io.Reader) (string) {
	buffer := new(s_bytes.Buffer)
	buffer.ReadFrom(reader)

	return buffer.String()
}

func IsExistPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ReadFileContent(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ObjectToJsonString(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if nil == err {
		return "", err
	}

	return string(bytes), nil
}

func Float64ToString(value float64) (string) {
	// -1 保留小数点几位
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func Float32ToString(value float32) (string) {
	return strconv.FormatFloat(float64(value), 'f', -1, 64)
}

func Int32ToString(value int) (string) {
	return strconv.Itoa(value)
}

func Int64ToString(value int64) (string) {
	return strconv.FormatInt(value, 10)
}
