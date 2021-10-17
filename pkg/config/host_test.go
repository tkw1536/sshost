package config_test

import (
	"reflect"
	"testing"

	"github.com/tkw1536/sshost/pkg/config"
)

func TestParseHost(t *testing.T) {
	tests := []struct {
		name    string
		wantH   config.Host
		wantErr bool
	}{
		{
			name:    "example.com",
			wantH:   config.Host{Host: "example.com"},
			wantErr: false,
		},
		{
			name:    "example.com:2222",
			wantH:   config.Host{Host: "example.com", Port: 2222},
			wantErr: false,
		},
		{
			name:    "user@example.com",
			wantH:   config.Host{Host: "example.com", User: "user"},
			wantErr: false,
		},
		{
			name:    "user@example.com:2222",
			wantH:   config.Host{Host: "example.com", User: "user", Port: 2222},
			wantErr: false,
		},
		{
			name:    "user@example.com:abcd",
			wantH:   config.Host{},
			wantErr: true,
		},
		{
			name:    "ssh://example.com",
			wantH:   config.Host{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotH, err := config.ParseHost(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotH, tt.wantH) {
				t.Errorf("ParseHost() = %v, want %v", gotH, tt.wantH)
			}
		})
	}
}
