package writer

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetOutput(io.Discard)
}

// TestNewTextWriter implenets unittest for NewTextWriter
func TestNewTextWriter(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		filename string
		header   string
		want     *TextWriter
	}{
		{
			name: "should return an instance of TextWriter",
			want: &TextWriter{
				path:     ".",
				filename: "test.js",
				header:   "// header",
			},
			path:     ".",
			filename: "test.js",
			header:   "// header",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextWriter(tt.path, tt.filename, tt.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTextWriter_Compile implements unit tests for TextWriter.Compile
func TestTextWriter_Compile(t *testing.T) {
	logrus.SetLevel(logrus.PanicLevel)
	data := struct {
		Text string
	}{
		Text: "World",
	}

	validTemplate := `Hello {{.Text}}`
	invalidDataTemplate := `Hello {{.Name}}`
	invalidFormatTemplate := `Hello {{.Text}`

	type fields struct {
		data []byte
	}
	type args struct {
		ctx  context.Context
		t    string
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    string
	}{

		{
			name: "should throw error for an invalid format template",
			args: args{
				t:    invalidFormatTemplate,
				data: data,
				ctx:  context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "should throw error for if there is data missmatch with template",
			args: args{
				t:    invalidDataTemplate,
				data: data,
				ctx:  context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "should generate code for a valid template",
			args: args{
				t:    validTemplate,
				data: data,
				ctx:  context.TODO(),
			},
			wantErr: false,
			want:    "Hello World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &TextWriter{
				data: tt.fields.data,
			}
			if err := w.Compile(tt.args.ctx, tt.args.t, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("TextWriter.Compile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.want != string(w.data) {
				t.Errorf("TextWriter.Compile(), TextWriter.data= %v, want %v", w.data, tt.want)
			}
		})
	}
}

// TestTextWriter_Format implements unit tests for TextWriter.Format
func TestTextWriter_Format(t *testing.T) {
	type fields struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		ctx     context.Context
	}{
		{
			name:    "should not throw any errors",
			wantErr: false,
			ctx:     context.TODO(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &TextWriter{
				data: tt.fields.data,
			}
			if err := w.Format(tt.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TextWriter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTextWriter_Write(t *testing.T) {
	type fields struct {
		data     []byte
		path     string
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should throw error if could not create path",
			fields: fields{
				data:     []byte("test"),
				path:     "/etc/files/",
				filename: "text.txt",
			},
			wantErr: true,
		},
		{
			name: "should throw error if count not open file",
			fields: fields{
				data:     []byte("test"),
				path:     "/",
				filename: "text.txt",
			},
			wantErr: true,
		},
		{
			name: "should return null if no errors",
			fields: fields{
				data:     []byte("test"),
				path:     "/tmp",
				filename: "text.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &TextWriter{
				data:     tt.fields.data,
				path:     tt.fields.path,
				filename: tt.fields.filename,
			}
			if err := w.Write(context.TODO()); (err != nil) != tt.wantErr {
				t.Errorf("TextWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
