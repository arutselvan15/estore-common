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
		logger = uLog.NewLogger().SetCluster(config.ClusterName).SetApplication(
			config.Application).SetComponent(resource).SetLevel(config.LogLevel)
	})

	return logger
}

// GetNewLogger returns a new log object
func GetNewLogger(resource string) uLog.CommonLog {
	return uLog.NewLogger().SetCluster(config.ClusterName).SetApplication(
		config.Application).SetComponent(resource).SetLevel(config.LogLevel)
}
