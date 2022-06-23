package writer

import (
	"context"
	"reflect"
	"testing"
)

func Test_getFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "should return file extension if present",
			filename: "main.go",
			want:     "go",
		},
		{
			name:     "should return file extension for complex names",
			filename: "server.main.go",
			want:     "go",
		},
		{
			name:     "should return expty string if no extension found",
			filename: "Makefile",
			want:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileExtension(tt.filename); got != tt.want {
				t.Errorf("getFileExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWriter(t *testing.T) {

	tests := []struct {
		name     string
		path     string
		filename string
		want     IWriter
	}{
		{
			name:     "should return go writer for go files",
			path:     "/tmp",
			filename: "main.go",
			want:     NewGoWriter("/tmp", "main.go"),
		},
		{
			name:     "should return text writer for all other files",
			path:     "/tmp",
			filename: "Makefile",
			want:     NewTextWriter("/tmp", "Makefile", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetWriter(tt.path, tt.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompileAndWrite(t *testing.T) {
	type args struct {
		path     string
		filename string
		tmpl     string
		data     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return error if failed to compile templates",
			args: args{
				path:     "/tmp",
				filename: "test.txt",
				tmpl:     "Hello {{.Nam}}",
				data: struct {
					Name string
				}{Name: "world"},
			},
			wantErr: true,
		},
		{
			name: "should return error if failed write file",
			args: args{
				path:     "/",
				filename: "test.txt",
				tmpl:     "Hello {{.Name}}",
				data: struct {
					Name string
				}{Name: "world"},
			},
			wantErr: true,
		},
		{
			name: "should return nil if generated successfully",
			args: args{
				path:     "/tmp",
				filename: "test.txt",
				tmpl:     "Hello {{.Name}}",
				data: struct {
					Name string
				}{Name: "world"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := CompileAndWrite(context.Background(), tt.args.path, tt.args.filename, tt.args.tmpl, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("CompileAndWrite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
