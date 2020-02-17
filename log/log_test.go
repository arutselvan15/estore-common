package log

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/arutselvan15/estore-common/config"
)

func TestLogger(t *testing.T) {
	_ = config.LoadFixture(config.FixtureDir)

	l := GetLogger("test")

	if l == nil {
		t.Error("GetLogger() is nil, want logger")
	}

	l = GetNewLogger("test")

	if l == nil {
		t.Error("GetNewLogger() is nil, want logger")
	}

	l = getLogger("test")

	if l == nil {
		t.Error("getLogger() is nil, want logger")
	}

	viper.Set("app.log.file.enabled", true)

	l = getLogger("test")

	if l == nil {
		t.Error("getLogger() file is nil, want logger")
	}
}
