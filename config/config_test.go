package config

import "testing"

func TestLoadFixture(t *testing.T) {
	type args struct {
		dir string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success Loadfixture working",
			args:    args{dir: "../fixture"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadFixture(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("LoadFixture() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
