package telemetry

import (
	"context"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name        string
		endpoint    string
		wantEnabled bool
		wantErr     bool
	}{
		{
			name:        "正常系: endpoint未指定ならトレーシングを無効化する",
			endpoint:    "",
			wantEnabled: false,
			wantErr:     false,
		},
		{
			name:        "正常系: endpoint指定時にトレーシングを初期化できる",
			endpoint:    "http://localhost:4318",
			wantEnabled: true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shutdown, enabled, err := Init(context.Background(), tt.endpoint, ServiceNameAPI, "dev")

			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if enabled != tt.wantEnabled {
				t.Errorf("Init() enabled = %v, want %v", enabled, tt.wantEnabled)
			}

			if err == nil {
				if shutdown == nil {
					t.Fatal("Init() shutdown is nil")
				}
				if err := shutdown(context.Background()); err != nil {
					t.Errorf("shutdown() error = %v", err)
				}
			}
		})
	}
}
