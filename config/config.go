// Package config provides common configurations
package config

import (
	"fmt"
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
)

// LogConfig lc
type LogConfig struct {
	Level          gLog.LevelLog
	Format         gLog.FormatterType
	LogFileEnabled bool
	LogFilePath    string
	LogFileFormat  gLog.FormatterType
	LogFileSize    int
	LogFileAge     int
	LogFileBkup    int
}

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

	_ = viper.BindEnv("app.name", "APP_NAME")
	_ = viper.BindEnv("app.freeze.startTime", "FREEZE_START_TIME")
	_ = viper.BindEnv("app.freeze.endTime", "FREEZE_END_TIME")
	_ = viper.BindEnv("app.freeze.message", "FREEZE_MESSAGE")
	_ = viper.BindEnv("app.freeze.components", "FREEZE_COMPONENTS")

	_ = viper.BindEnv("app.system.namespaces", "SYSTEM_NAMESPACES")
	_ = viper.BindEnv("app.system.users", "SYSTEM_USERS")

	_ = viper.BindEnv("app.blacklist.namespaces", "BLACKLIST_NAMESPACES")
	_ = viper.BindEnv("app.blacklist.users", "BLACKLIST_USERS")

	_ = viper.BindEnv("app.log.level", "LOG_LEVEL")
	_ = viper.BindEnv("app.log.format", "LOG_FORMAT")
	_ = viper.BindEnv("app.log.file.enabled", "LOG_FILE_ENABLED")
	_ = viper.BindEnv("app.log.file.format", "LOG_FILE_FORMAT")
	_ = viper.BindEnv("app.log.file.dir", "LOG_FILE_DIR")
	_ = viper.BindEnv("app.log.file.name", "LOG_FILE_NAME")
	_ = viper.BindEnv("app.log.file.size", "LOG_FILE_SIZE")
	_ = viper.BindEnv("app.log.file.age", "LOG_FILE_AGE")
	_ = viper.BindEnv("app.log.file.backup", "LOG_FILE_BACKUP")

	_ = viper.BindEnv("cluster.name", "CLUSTER_NAME")
	_ = viper.BindEnv("cluster.kubeconfig", "KUBECONFIG")
}

// LoadFixture load fixtures
func LoadFixture(dir string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)

	return viper.ReadInConfig()
}

// GetAppName app name
func GetAppName() string {
	return viper.GetString("app.name")
}

// GetKubeConfigPath kube config
func GetKubeConfigPath() string {
	return viper.GetString("cluster.kubeconfig")
}

// GetLogConfig log config
func GetLogConfig() LogConfig {
	lc := LogConfig{}

	lc.Format = gLog.TextFormatterType
	if viper.GetString("app.log.format") == "json" {
		lc.Format = gLog.JSONFormatterType
	}

	lc.Level = gLog.InfoLevel
	if viper.GetString("app.log.level") == "debug" {
		lc.Level = gLog.DebugLevel
	}

	lc.LogFileEnabled = viper.GetBool("app.log.file.enabled")
	lc.LogFilePath = fmt.Sprintf("%s/%s", viper.GetString("app.log.file.dir"), viper.GetString("app.log.file.name"))
	lc.LogFileAge = viper.GetInt("app.log.file.age")
	lc.LogFileBkup = viper.GetInt("app.log.file.backup")
	lc.LogFileSize = viper.GetInt("app.log.file.size")

	lc.LogFileFormat = gLog.TextFormatterType
	if viper.GetString("app.log.format") == "json" {
		lc.LogFileFormat = gLog.JSONFormatterType
	}

	return lc
}

// GetClusterName cluster name
func GetClusterName() string {
	clusterName := viper.GetString("cluster.name")
	if clusterName == "" {
		clusterName = "unknown"
	}

	return clusterName
}

// GetSystemNamespaces sys ns
func GetSystemNamespaces() map[string]bool {
	return stringToBoolMap(viper.GetString("app.system.namespaces"), ",")
}

// GetSystemUsers sys users
func GetSystemUsers() map[string]bool {
	return stringToBoolMap(viper.GetString("app.system.users"), ",")
}

// GetBlacklistNamespaces black list ns
func GetBlacklistNamespaces() map[string]bool {
	return stringToBoolMap(viper.GetString("app.blacklist.namespaces"), ",")
}

// GetBlacklistUsers black list users
func GetBlacklistUsers() map[string]bool {
	return stringToBoolMap(viper.GetString("app.blacklist.users"), ",")
}

func stringToBoolMap(str, sep string) map[string]bool {
	strBoolMap := make(map[string]bool)

	if str != "" {
		for _, i := range strings.Split(str, sep) {
			strBoolMap[i] = true
		}
	}

	return strBoolMap
}
