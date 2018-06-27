package common

import (
	"log"
	"os"
	"time"
	"fmt"
	"sync"
)

var (
	defaultLogging *logging = nil
	loggingLock    sync.Mutex
)

func GetLogging() (*logging) {
	if nil != defaultLogging {
		if time.Now().Day() == defaultLogging.Day {
			return defaultLogging
		}
	}

	loggingLock.Lock()
	defer loggingLock.Unlock()

	if nil != defaultLogging {
		if time.Now().Day() == defaultLogging.Day {
			return defaultLogging
		}

		defaultLogging.Close()
	}

	defaultLogging = newLogging(false)

	return defaultLogging
}

func newLogging(toFile bool) (*logging) {
	logging := &logging{}

	now := time.Now()
	logging.Day = now.Day()
	logging.ok = false
	if !toFile {
		return logging
	}

	logPath := fmt.Sprintf("./goapp_%s.log", now.Format("20060102"))
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	logging.ok = nil == err

	if logging.ok {
		logging.logFile = logFile
		logging.logger = log.New(logFile, "", log.Ldate|log.Ltime)
	}

	return logging
}

type logging struct {
	logger  *log.Logger
	ok      bool
	logFile *os.File
	toFile  bool
	Day     int
}

func (myself *logging) Close() (error) {
	if myself.ok {
		return myself.logFile.Close()
	}

	return nil
}

func (myself *logging) Debug(err interface{}, format interface{}, parameters ...interface{}) {
	myself.log("Debug", err, format, parameters...)
}

func (myself *logging) Info(err interface{}, format interface{}, parameters ...interface{}) {
	myself.log("Info", err, format, parameters...)
}

func (myself *logging) Warn(err interface{}, format interface{}, parameters ...interface{}) {
	myself.log("Warn", err, format, parameters...)
}

func (myself *logging) Error(err interface{}, format interface{}, parameters ...interface{}) {
	myself.log("Error", err, format, parameters...)
}

func (myself *logging) Fatal(err interface{}, format interface{}, parameters ...interface{}) {
	myself.log("Fatal", err, format, parameters...)
}

func (myself *logging) log(level string, err interface{}, format interface{}, parameters ...interface{}) {
	defer func() {
		value := recover()
		if nil != value {
			log.Printf("failed to log message; error: %s", fmt.Sprint(value))
		}
	}()

	formatMessage := ""
	if nil == err {
		formatMessage = fmt.Sprintf("[%s][%s]", level, fmt.Sprint(format))
	} else {
		formatMessage = fmt.Sprintf("[%s][%s][error: %s]", level, fmt.Sprint(format), fmt.Sprint(err))
	}

	if myself.ok {
		myself.logger.Printf(formatMessage, parameters...)
	} else {
		log.Printf(formatMessage, parameters...)
	}
}
