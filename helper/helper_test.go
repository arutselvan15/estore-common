package helper

import (
	"reflect"
	"testing"
)

func TestContainsString(t *testing.T) {
	type args struct {
		slice []string
		s     string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "success container string found", args: args{slice: []string{"test"}, s: "test"}, want: true},
		{name: "success container string not found", args: args{slice: []string{"test1"}, s: "test"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsString(tt.args.slice, tt.args.s); got != tt.want {
				t.Errorf("ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveString(t *testing.T) {
	type args struct {
		slice []string
		s     string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "success remove string", args: args{slice: []string{"test", "test1"}, s: "test"}, want: []string{"test1"}},
		{name: "success remove string not found", args: args{slice: []string{"test1"}, s: "test2"}, want: []string{"test1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveString(tt.args.slice, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveString() = %v, want %v", got, tt.want)
			}
		})
	}
}
