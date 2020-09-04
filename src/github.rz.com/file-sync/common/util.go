package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	random      = rand.New(rand.NewSource(time.Now().UnixNano()))
	defaultIpV4 = ""
)

func IsStringBlank(value string) bool {
	return 0 == len(strings.TrimSpace(value))
}

func ReaderToString(reader io.Reader) string {
	if nil == reader {
		return ""
	}

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(reader)
	if nil != err {
		fmt.Println("error from ReaderToString; err", err)
	}

	return buffer.String()
}

func IsExistPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ReadFileContent(filePath string) (string, error) {
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func ObjectToJsonString(value interface{}) string {
	buffer, err := json.Marshal(value)
	if nil != err {
		return ""
	}

	return string(buffer)
}

func Float64ToString(value float64) string {
	// -1 保留小数点几位
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func Float32ToString(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 64)
}

func Int32ToString(value int) string {
	return strconv.Itoa(value)
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func StringToInt32(value string) (int, error) {
	return strconv.Atoi(value)
}

func StringToInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func StringToFloat32(value string) (float32, error) {
	convertedValue, err := strconv.ParseFloat(value, 32)
	if nil != err {
		return 0, err
	}
	return float32(convertedValue), nil
}

func StringToFloat64(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// 字符串截取
// 区间 [start, end) start 包含  end 不含
// 支持中文
func Substring(source string, start int, end int) (string, error) {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return "", errors.New(fmt.Sprintf("The [substring] start(%d) < 0 || end(%d) > length(%d) || start(%d) > end(%d)",
			start, end, length, start, end))
	}

	if start == 0 && end == length {
		return source, nil
	}

	var result = ""
	for i := start; i < end; i++ {
		result += string(r[i])
	}

	return result, nil
}

func SortInt64(values []int64) {
	sort.Slice(values, func(currentIndex int, nextIndex int) bool { return values[currentIndex] < values[nextIndex] })
}

func SortReverseInt64(values []int64) {
	sort.Slice(values, func(currentIndex int, nextIndex int) bool { return values[currentIndex] > values[nextIndex] })
}

func SortInt(values []int) {
	sort.Slice(values, func(currentIndex int, nextIndex int) bool { return values[currentIndex] < values[nextIndex] })
}

func SortReverseInt(values []int) {
	sort.Slice(values, func(currentIndex int, nextIndex int) bool { return values[currentIndex] > values[nextIndex] })
}

func SortString(values []string) {
	sort.Strings(values)
}

func SortReverseString(values []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(values)))
}

func GetIpV4s() ([]string, error) {
	interfaceAddresses, err := net.InterfaceAddrs()
	if nil != err {
		return nil, err
	}

	var ipV4s []string
	for _, interfaceAddress := range interfaceAddresses {
		ipNet, ok := interfaceAddress.(*net.IPNet)
		if ok && nil != ipNet.IP && !ipNet.IP.IsLoopback() && nil != ipNet.IP.To4() {
			ipV4s = append(ipV4s, ipNet.IP.String())
		}
	}

	return ipV4s, nil
}

func GetDefaultIpV4() string {
	if "" != defaultIpV4 {
		return defaultIpV4
	}

	ipV4s, err := GetIpV4s()
	if nil != err || 0 == len(ipV4s) {
		defaultIpV4 = "127.0.0.1"
	}
	defaultIpV4 = ipV4s[0]

	return defaultIpV4
}

func MD5(value []byte) string {
	hash := md5.New()
	hash.Write(value)
	buffer := hash.Sum(nil)

	return hex.EncodeToString(buffer)
}

func RandomStringV1(length int) string {
	if length < 1 {
		length = 1
	}

	buffer := make([]byte, length)
	for i := 0; i < length; i++ {
		value := random.Intn(26)
		if 0 == value%3 {
			value = value + 65
		} else if 1 == value%3 {
			value = value + 97
		} else {
			value = value%10 + 48
		}

		buffer[i] = byte(value)
	}

	return string(buffer)
}

func StringHash(value string) int {
	hash := int(crc32.ChecksumIEEE([]byte(value)))
	if hash >= 0 {
		return hash
	}
	if -hash >= 0 {
		return -hash
	}
	return 0
}
