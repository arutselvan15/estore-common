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
	// LogFormat log level
	LogFormat gLog.FormatterType
	// KubeConfigPath kube config path
	KubeConfigPath string
	// WhitelistUsers white list users
	WhitelistUsers []string
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
	_ = viper.BindEnv("app.whitelist.users", "WHITELIST_USERS")
	_ = viper.BindEnv("app.log.level", "LOG_LEVEL")
	_ = viper.BindEnv("app.log.format", "LOG_FORMAT")
	_ = viper.BindEnv("app.log.file.enabled", "LOG_FILE_ROTATE")
	_ = viper.BindEnv("app.log.file.format", "LOG_FILE_FORMAT")
	_ = viper.BindEnv("app.log.file.dir", "LOG_FILE_DIR")
	_ = viper.BindEnv("app.log.file.name", "LOG_FILE_NAME")
	_ = viper.BindEnv("app.log.file.size", "LOG_FILE_SIZE")
	_ = viper.BindEnv("app.log.file.age", "LOG_FILE_AGE")
	_ = viper.BindEnv("app.log.file.backup", "LOG_FILE_BACKUP")

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

	LogFormat = gLog.TextFormatterType
	if viper.GetString("app.log.format") == "json" {
		LogFormat = gLog.JSONFormatterType
	}

	if viper.GetString("app.whitelist.namespaces") != "" {
		WhitelistNamespaces = strings.Split(viper.GetString("app.whitelist.namespaces"), ",")
	}

	if viper.GetString("app.whitelist.namespaces") != "" {
		WhitelistNamespaces = strings.Split(viper.GetString("app.whitelist.users"), ",")
	}

	KubeConfigPath = viper.GetString("cluster.kubeconfig")
}

// LoadFixture load fixtures
func LoadFixture(dir string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)

	return viper.ReadInConfig()
}
