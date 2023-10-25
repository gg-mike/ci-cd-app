package sys

import (
	"testing"
)

func TestGetEnvWithFallback(t *testing.T) {
	type fields struct {
		key 		 string
		fallback string
	}
	tests := []struct {
		name  string
		args  fields
		isEnv bool
		want  string
	}{
		{ "Env present", fields{ "FOO", "bar" }, true, "foo" },
		{ "Env absent", fields{ "FOO", "bar" }, false, "bar" },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isEnv {
				t.Setenv(tt.args.key, tt.want)
			}
			if got := GetEnvWithFallback(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnvWithFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequiredEnv(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		isEnv 	bool
		want    string
		wantErr bool
	}{
		{ "Env present", "FOO", true, "foo", false },
		{ "Env absent", "FOO", false, "", true },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isEnv {
				t.Setenv(tt.args, tt.want)
			}
			got, err := GetRequiredEnv(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRequiredEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRequiredEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
