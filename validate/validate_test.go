package validate

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

func TestAdmissionRequired(t *testing.T) {
	type args struct {
		admissionAnnotationKey string
		metadata               *metav1.ObjectMeta
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success admission required", want: true,
			args: args{"", &metav1.ObjectMeta{Namespace: ""}},
		},
		{
			name: "success admission not required special annotation", want: false,
			args: args{"validate", &metav1.ObjectMeta{Annotations: map[string]string{"validate": "no"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := AdmissionRequired(tt.args.admissionAnnotationKey, tt.args.metadata); got != tt.want {
				t.Errorf("AdmissionRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreatePatchAnnotations(t *testing.T) {
	type args struct {
		availableAnnotations map[string]string
		addAnnotations       map[string]string
	}

	tests := []struct {
		name      string
		args      args
		wantPatch []PatchOperation
	}{
		{
			name:      "success new patch with no current annotations",
			args:      args{availableAnnotations: nil, addAnnotations: map[string]string{"key1": "val1"}},
			wantPatch: []PatchOperation{{Op: "add", Path: "/metadata/annotations", Value: map[string]string{"key1": "val1"}}},
		},
		{
			name:      "success add patch with current annotations",
			args:      args{availableAnnotations: map[string]string{}, addAnnotations: map[string]string{"key1": "val1"}},
			wantPatch: []PatchOperation{{Op: "add", Path: "/metadata/annotations/key1", Value: "val1"}},
		},
		{
			name:      "success replace patch with current annotations",
			args:      args{availableAnnotations: map[string]string{"key1": "val0"}, addAnnotations: map[string]string{"key1": "val1"}},
			wantPatch: []PatchOperation{{Op: "replace", Path: "/metadata/annotations/key1", Value: "val1"}},
		},
		{
			name:      "success add with slash annotations",
			args:      args{availableAnnotations: nil, addAnnotations: map[string]string{"resource.app.com/key1": "val1"}},
			wantPatch: []PatchOperation{{Op: "add", Path: "/metadata/annotations", Value: map[string]string{"resource.app.com~1key1": "val1"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPatch := CreatePatchAnnotations(tt.args.availableAnnotations, tt.args.addAnnotations); !reflect.DeepEqual(gotPatch, tt.wantPatch) {
				t.Errorf("CreatePatchAnnotations() = %v, want %v", gotPatch, tt.wantPatch)
			}
		})
	}
}

func TestCreatePatchLabels(t *testing.T) {
	type args struct {
		availableLabels map[string]string
		addLabels       map[string]string
	}

	tests := []struct {
		name      string
		args      args
		wantPatch []PatchOperation
	}{
		{
			name:      "success new patch with no current labels",
			args:      args{availableLabels: nil, addLabels: map[string]string{"key1": "val1"}},
			wantPatch: []PatchOperation{{Op: "add", Path: "/metadata/labels", Value: map[string]string{"key1": "val1"}}},
		},
		{
			name:      "success add patch with current labels",
			args:      args{availableLabels: map[string]string{}, addLabels: map[string]string{"key1": "val1"}},
			wantPatch: []PatchOperation{{Op: "add", Path: "/metadata/labels/key1", Value: "val1"}},
		},
		{
			name:      "success replace patch with current labels",
			args:      args{availableLabels: map[string]string{"key1": "val0"}, addLabels: map[string]string{"key1": "val1"}},
			wantPatch: []PatchOperation{{Op: "replace", Path: "/metadata/labels/key1", Value: "val1"}},
		},
		{
			name:      "success add with slash",
			args:      args{availableLabels: nil, addLabels: map[string]string{"resource.app.com/key1": "val1"}},
			wantPatch: []PatchOperation{{Op: "add", Path: "/metadata/labels", Value: map[string]string{"resource.app.com~1key1": "val1"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPatch := CreatePatchLabels(tt.args.availableLabels, tt.args.addLabels); !reflect.DeepEqual(gotPatch, tt.wantPatch) {
				t.Errorf("CreatePatchLabels() = %v, want %v", gotPatch, tt.wantPatch)
			}
		})
	}
}

func Test_patchHandleSlash(t *testing.T) {
	type args struct {
		data map[string]string
	}

	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "success handle slash", args: args{map[string]string{"hello/hello/hello": "hello/hello"}},
			want: map[string]string{"hello~1hello~1hello": "hello/hello"},
		},
		{
			name: "success handle normal", args: args{map[string]string{"hello+hello": "hello/hello", "hello:hello": "hello/hello"}},
			want: map[string]string{"hello+hello": "hello/hello", "hello:hello": "hello/hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := patchHandleSlash(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("patchHandleSlash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckBlacklistNamespace(t *testing.T) {
	type args struct {
		ns string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "success check black list ns", args: args{ns: "test"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckBlacklistNamespace(tt.args.ns); got != tt.want {
				t.Errorf("CheckBlacklistNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckBlacklistUser(t *testing.T) {
	type args struct {
		user string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "success check black list user", args: args{user: "test"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckBlacklistUser(tt.args.user); got != tt.want {
				t.Errorf("CheckBlacklistUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckSystemNamespace(t *testing.T) {
	type args struct {
		ns string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "success check system ns", args: args{ns: "test"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSystemNamespace(tt.args.ns); got != tt.want {
				t.Errorf("CheckSystemNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckSystemUser(t *testing.T) {
	type args struct {
		user string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "success check system user", args: args{user: "test"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSystemUser(tt.args.user); got != tt.want {
				t.Errorf("CheckSystemUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkMatch(t *testing.T) {
	type args struct {
		value   string
		options map[string]bool
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "success perfect match", args: args{value: "test", options: map[string]bool{"test": true}}, want: true},
		{name: "success prefix match", args: args{value: "test-user", options: map[string]bool{"test": true}}, want: true},
		{name: "success no match", args: args{value: "test", options: map[string]bool{"unknown": true}}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkMatch(tt.args.value, tt.args.options); got != tt.want {
				t.Errorf("checkMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
