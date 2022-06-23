package writer

import (
	"context"
	"os"
	"reflect"
	"strings"
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
					header:   autoGenText,
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

func TestGoWriter_Format(t *testing.T) {
	code := `package main
	func main() {
		fmt.Println("test")
	}
	`

	w := NewGoWriter("/tmp", "main.go")
	w.data = []byte(code)
	ctx := context.TODO()
	w.Write(ctx)
	w.Format(ctx)

	data, _ := os.ReadFile("/tmp/main.go")

	if !strings.Contains(string(data), "import \"fmt\"") {
		t.Errorf("code did not formatted as expected, check go-imports version")
	}
}
