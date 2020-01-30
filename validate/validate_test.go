package validate

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/arutselvan15/estore-common/config"
)

func Test_checkFreezeEnabled(t *testing.T) {
	_ = config.LoadFixture(config.FixtureDir)

	const timeToIncrease = 2

	timeNow := time.Now()

	type args struct {
		startTime string
		endTime   string
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "success - freeze enabled",
			args:    args{startTime: timeNow.Add(-5 * time.Hour).Format(config.TimeLayout), endTime: timeNow.Add(timeToIncrease * time.Hour).Format(config.TimeLayout)},
			want:    true,
			wantErr: false,
		},
		{
			name:    "success - long freeze",
			args:    args{startTime: timeNow.Add(-5 * time.Hour).Format(config.TimeLayout), endTime: ""},
			want:    true,
			wantErr: false,
		},
		{
			name:    "success - no freeze data",
			args:    args{startTime: "", endTime: ""},
			want:    false,
			wantErr: false,
		},
		{
			name:    "success - freeze over",
			args:    args{startTime: timeNow.Add(-5 * time.Hour).Format(config.TimeLayout), endTime: timeNow.Add(-3 * time.Hour).Format(config.TimeLayout)},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkFreezeEnabled(tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkFreezeEnabled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkFreezeEnabled() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFreezeEnabled(t *testing.T) {
	_ = config.LoadFixture(config.FixtureDir)

	// restore fixture
	restoreComponents := viper.GetString("app.freeze.components")
	defer func(restoreComponents string) {
		viper.Set("app.freeze.components", restoreComponents)
	}(restoreComponents)

	type args struct {
		component            string
		mockFreezeEnabled    bool
		mockErr              error
		monkFreezeComponents string
	}

	tests := []struct {
		name, want1   string
		args          args
		want, wantErr bool
	}{
		{
			name: "success no freeze", args: args{component: "", monkFreezeComponents: "all", mockFreezeEnabled: false}, want: false, want1: "", wantErr: false,
		},
		{
			name: "success freeze for all", args: args{component: "", monkFreezeComponents: "all", mockFreezeEnabled: true}, want: true, want1: viper.GetString("app.freeze.message"), wantErr: false,
		},
		{
			name: "success freeze but component not involved", args: args{component: "dummy1", monkFreezeComponents: "unknown", mockFreezeEnabled: true}, want: false, want1: "", wantErr: false,
		},
		{
			name: "success freeze enabled for component", args: args{component: "dummy", monkFreezeComponents: "dummy", mockFreezeEnabled: true}, want: true, want1: viper.GetString("app.freeze.message"), wantErr: false,
		},
		{
			name: "failure error during freeze check", args: args{component: "dummy", monkFreezeComponents: "dummy", mockFreezeEnabled: false, mockErr: fmt.Errorf("")}, want: false, want1: "", wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock component values - restore using defer to avoid changes in fixture
			viper.Set("app.freeze.components", tt.args.monkFreezeComponents)

			checkFreezeEnabled = func(startTime, endTime string) (b bool, err error) {
				return tt.args.mockFreezeEnabled, tt.args.mockErr
			}

			got, got1, err := FreezeEnabled(tt.args.component)
			if (err != nil) != tt.wantErr {
				t.Errorf("FreezeEnabled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FreezeEnabled() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FreezeEnabled() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getFreezeComponents(t *testing.T) {
	_ = config.LoadFixture(config.FixtureDir)

	tests := []struct {
		name string
		want map[string]bool
	}{
		{
			name: "success get components",
			want: map[string]bool{"all": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFreezeComponents(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFreezeComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}
