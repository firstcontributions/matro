package cleaner

import (
	"os"
	"testing"
)

func TestClean(t *testing.T) {
	type args struct {
		basePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "should clean the dirs",
			args:    args{basePath: "__test"},
			wantErr: false,
		},
		{
			name:    "should clean the dirs",
			args:    args{basePath: "__test_fail/test.txt"},
			wantErr: true,
		},
	}
	os.MkdirAll("__test_fail", 0777)
	f, _ := os.Create("__test_fail/test.txt")
	f.WriteString("test file")
	f.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Clean(tt.args.basePath); (err != nil) != tt.wantErr {
				t.Errorf("Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.RemoveAll("__test_fail")
}
