// Package config provides common configurations
package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	// FixtureDir fixture diretory
	FixtureDir = "../fixture"
	// TimeLayout time layout
	TimeLayout = time.RFC3339
	// Component component name
	Component = "estore"
	// SubComponent sub component
	SubComponent = "common"
)

var (
	// ClusterName cluster name
	ClusterName string
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
	_ = viper.BindEnv("cluster.name", "CLUSTER_NAME")

	ClusterName = viper.GetString("cluster.name")
	if ClusterName == "" {
		ClusterName = "unknown"
	}
}

// LoadFixture load fixtures
func LoadFixture(dir string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)

	return viper.ReadInConfig()
}
