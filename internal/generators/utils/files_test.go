package utils

import (
	"os"
	"testing"
)

func TestEnsurePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "should throw error if cannot create directory",
			path:    "/usr/bin/test",
			wantErr: true,
		},
		{
			name:    "should not raise errors if path already exists",
			path:    "/tmp",
			wantErr: false,
		},
		{
			name:    "should not raise errors if path can be created",
			path:    "/tmp/go-test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EnsurePath(tt.path); (err != nil) != tt.wantErr {
				t.Errorf("EnsurePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	os.RemoveAll("/tmp/go-test")
}
