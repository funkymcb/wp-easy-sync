package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic",
			args:    args{"./fixtures/basic.yaml"},
			wantErr: false,
		},
		{
			name:    "Corrupted",
			args:    args{"./fixtures/corrupted.yaml"},
			wantErr: true,
		},
		{
			name:    "Not Found",
			args:    args{"./fixtures/file_does_not_exist.yaml"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConfig(tt.args.configPath); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "Basic",
			args: args{"./fixtures/basic.yaml"},
			want: &Config{
				Easyverein: EasyvereinCfg{
					Host:  "easyverein.com",
					Path:  "/api/stable/",
					Token: "some nice token",
					Options: map[string]string{
						"limit": "100",
					},
				},
				Wordpress: WordpressCfg{
					Host:     "your-domain.com",
					Path:     "/wp-json/wp/v2/",
					Username: "username",
					Password: "password",
					Options: map[string]string{
						"per_page": "100",
						"context":  "edit",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConfig(tt.args.configPath); err != nil {
				t.Errorf("LoadConfig() failed with %v", err)
			}

			if got := GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
