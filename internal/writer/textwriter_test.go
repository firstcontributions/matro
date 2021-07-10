package writer

import (
	"context"
	"reflect"
	"testing"
)

// TestNewTextWriter implenets unittest for NewTextWriter
func TestNewTextWriter(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		filename string
		want     *TextWriter
	}{
		{
			name: "should return an instance of TextWriter",
			want: &TextWriter{
				path:     ".",
				filename: "test.js",
			},
			path:     ".",
			filename: "test.js",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextWriter(tt.path, tt.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTextWriter_Compile implements unit tests for TextWriter.Compile
func TestTextWriter_Compile(t *testing.T) {
	data := struct {
		Text string
	}{
		Text: "World",
	}

	validTemplate := `Hello {{.Text}}`
	inValidTemplate := `Hello {{.Name}}`

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
			name: "should generate code for a valid template",
			args: args{
				t:    validTemplate,
				data: data,
				ctx:  context.TODO(),
			},
			wantErr: false,
			want:    "Hello World",
		},
		{
			name: "should throw error for an invalid template",
			args: args{
				t:    inValidTemplate,
				data: data,
				ctx:  context.TODO(),
			},
			wantErr: true,
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
