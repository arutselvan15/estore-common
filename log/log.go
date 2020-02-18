// Package log provides common logs
package log

import (
	"sync"

	uLog "github.com/arutselvan15/go-utils/log"

	"github.com/arutselvan15/estore-common/config"
)

var (
	once   sync.Once
	logger uLog.CommonLog
)

// GetLogger returns a singleton of the Log object
func GetLogger(resource string) uLog.CommonLog {
	once.Do(func() {
		logger = getLogger(resource)
	})

	return logger
}

// GetNewLogger returns a new log object
func GetNewLogger(resource string) uLog.CommonLog {
	return getLogger(resource)
}

func getLogger(resource string) uLog.CommonLog {
	var newLogger uLog.CommonLog

	lc := config.GetLogConfig()

	if lc.LogFileEnabled {
		newLogger = uLog.NewLoggerWithFile(lc.LogFilePath, lc.LogFileSize, lc.LogFileAge, lc.LogFileBkup)
		newLogger.SetLogFileFormatterType(lc.LogFileFormat)
	} else {
		newLogger = uLog.NewLogger()
	}

	newLogger.SetCluster(config.GetClusterName()).SetApplication(
		config.GetAppName()).SetResource(resource).SetLevel(
		lc.Level).SetFormatterType(lc.Format)

	return newLogger
}
