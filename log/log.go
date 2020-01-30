// Package log provides common logs
package log

import (
	"sync"

	gLog "github.com/arutselvan15/go-utils/log"

	"github.com/arutselvan15/estore-common/config"
)

var (
	once        sync.Once
	logInstance gLog.CustomLog
)

// GetInstance returns a singleton of the Log object
func GetInstance() gLog.CustomLog {
	once.Do(func() {
		logInstance = gLog.NewLogger().SetComponent(config.Component).SetSubComponent(config.SubComponent).SetCluster(config.ClusterName)
	})

	return logInstance
}
