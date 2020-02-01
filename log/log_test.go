package log

import (
	"testing"

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
}
