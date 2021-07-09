package writer

import (
	"reflect"
	"testing"
)

// TestNewGoWriter tests the constructor for GoWriter
func TestNewGoWriter(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		filename string
		want     *GoWriter
	}{
		{
			name:     "should construct an instance of GoWriter",
			path:     ".",
			filename: "test.js",
			want: &GoWriter{
				TextWriter: &TextWriter{
					filename: "test.js",
					path:     ".",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGoWriter(tt.path, tt.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGoWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}
