// Package log provides common logs
package log

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"

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

	// log file config
	if viper.GetBool("app.log.file.enabled") {
		fName := fmt.Sprintf("%s/%s", viper.GetString("app.log.file.dir"),
			viper.GetString("app.log.file.name"))

		newLogger = uLog.NewLoggerWithFile(fName, viper.GetInt("app.log.file.size"),
			viper.GetInt("app.log.file.age"), viper.GetInt("app.log.file.backup"))

		switch strings.ToLower(viper.GetString("app.log.file.format")) {
		case "text":
			newLogger.SetLogFileFormatterType(uLog.TextFormatterType)
		case "json":
			newLogger.SetLogFileFormatterType(uLog.JSONFormatterType)
		default:
			newLogger.SetLogFileFormatterType(uLog.TextFormatterType)
		}
	} else {
		newLogger = uLog.NewLogger()
	}

	newLogger.SetCluster(config.ClusterName).SetApplication(
		config.Application).SetResource(resource).SetLevel(
		config.LogLevel).SetFormatterType(config.LogFormat)

	return newLogger
}
