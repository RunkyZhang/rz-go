package common

import (
	"strings"
	"bytes"
	"io"
)

func IsStringBlank(value string) (bool) {
	return 0 == len(strings.TrimSpace(value))
}

func ReaderToString(reader io.Reader) (string) {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)

	return buffer.String()
}
