package common

import (
	"strings"
	s_bytes "bytes"
	"io"
	"os"
	"io/ioutil"
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
