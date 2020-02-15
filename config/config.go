// Package config provides common configurations
package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"

	gLog "github.com/arutselvan15/go-utils/log"
)

const (
	// FixtureDir fixture directory
	FixtureDir = "../fixture"
	// TimeLayout time layout
	TimeLayout = time.RFC3339
	// Application app name
	Application = "estore"
)

var (
	// ClusterName cluster name
	ClusterName string
	// LogLevel log level
	LogLevel gLog.LevelLog
	// KubeConfigPath kube config path
	KubeConfigPath string
	// WhitelistNamespaces whitelist namespace
	WhitelistNamespaces []string
)

func init() {
	defaultConfigPaths := []string{
		"/etc/viper",
		"./",
	}

	viper.SetConfigName("config")

	for _, p := range defaultConfigPaths {
		viper.AddConfigPath(p)
	}

	// errors ignored
	_ = viper.ReadInConfig()

	_ = viper.BindEnv("app.freeze.startTime", "FREEZE_START_TIME")
	_ = viper.BindEnv("app.freeze.endTime", "FREEZE_END_TIME")
	_ = viper.BindEnv("app.freeze.message", "FREEZE_MESSAGE")
	_ = viper.BindEnv("app.freeze.components", "FREEZE_COMPONENTS")
	_ = viper.BindEnv("app.whitelist.namespaces", "WHITELIST_NAMESPACES")
	_ = viper.BindEnv("app.log.level", "LOG_LEVEL")

	_ = viper.BindEnv("cluster.name", "CLUSTER_NAME")
	_ = viper.BindEnv("cluster.kubeconfig", "KUBECONFIG")

	ClusterName = viper.GetString("cluster.name")
	if ClusterName == "" {
		ClusterName = "unknown"
	}

	LogLevel = gLog.InfoLevel
	if viper.GetString("app.log.level") == "debug" {
		LogLevel = gLog.DebugLevel
	}

	if viper.GetString("app.whitelist.namespaces") != "" {
		WhitelistNamespaces = strings.Split(viper.GetString("app.whitelist.namespaces"), ",")
	}

	KubeConfigPath = viper.GetString("cluster.kubeconfig")
}

// LoadFixture load fixtures
func LoadFixture(dir string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)

	return viper.ReadInConfig()
}
