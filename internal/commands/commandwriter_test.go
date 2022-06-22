package commands

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestNewCommandWriter(t *testing.T) {
	tests := []struct {
		name string
		args io.Writer
		want *CommandWriter
	}{
		{
			name: "it should return an instance of commandwriter",
			args: os.Stdout,
			want: &CommandWriter{writer: os.Stdout},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommandWriter(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommandWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandWriter_Write(t *testing.T) {
	type fields struct {
		writer io.Writer
	}
	tests := []struct {
		name        string
		fields      fields
		args        []string
		wantErr     bool
		wantToWrite string
	}{
		{
			name: "should write given strings to writer interface",
			fields: fields{
				writer: new(bytes.Buffer),
			},
			args:        []string{"hello matro", "test"},
			wantErr:     false,
			wantToWrite: "hello matro test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandWriter{
				writer: tt.fields.writer,
			}
			if err := c.Write(tt.args...); (err != nil) != tt.wantErr {
				t.Errorf("CommandWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if res := tt.fields.writer.(*bytes.Buffer).String(); res != tt.wantToWrite {
				t.Errorf("CommandWriter.Write() exptected to write  %v, but wrote %v", res, tt.wantToWrite)
			}
		})
	}
}
