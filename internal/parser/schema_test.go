package parser

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewDefinition(t *testing.T) {
	tests := []struct {
		name string
		want *Definition
	}{
		{
			name: "should create a new instance of Definition",
			want: &Definition{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefinition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestDefinition_ParseFrom(t *testing.T) {
	type fields struct {
		Modules  []Module
		Repo     string
		Queries  []*Type
		Defaults *Defaults
	}

	tests := []struct {
		name    string
		fields  fields
		args    io.Reader
		want    *Definition
		wantErr bool
	}{
		{
			name:    "should return error if there is any while reading the config",
			args:    errReader{},
			wantErr: true,
			want:    NewDefinition(),
		},
		{
			name:    "should return error if there is any while parsing the config",
			args:    strings.NewReader("{"),
			wantErr: true,
			want:    NewDefinition(),
		},
		{
			name:    "should parse given config",
			args:    strings.NewReader(`{"repo": "test"}`),
			wantErr: false,
			want:    &Definition{Repo: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				Modules:  tt.fields.Modules,
				Repo:     tt.fields.Repo,
				Queries:  tt.fields.Queries,
				Defaults: tt.fields.Defaults,
			}
			got, err := d.ParseFrom(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Definition.ParseFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Definition.ParseFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
