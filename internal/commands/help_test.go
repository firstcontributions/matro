package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelp_Help(t *testing.T) {
	txt := `
	matro help
	Help command prints help

	`
	type fields struct {
		CommandWriter *CommandWriter
	}
	tests := []struct {
		name        string
		fields      fields
		wantToWrite string
	}{
		{
			name: "should print help text",
			fields: fields{
				CommandWriter: NewCommandWriter(new(bytes.Buffer)),
			},
			wantToWrite: txt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Help{
				CommandWriter: tt.fields.CommandWriter,
			}
			h.Help()
			if res := h.writer.(*bytes.Buffer).String(); string(res) != string(tt.wantToWrite) {
				t.Errorf("CommandWriter.Write() exptected to write  %v, but wrote %v", string(res), string(tt.wantToWrite))
			}
		})
	}
}

func TestHelp_Exec(t *testing.T) {
	expectedOutput := `Usage
            
	matro help
	Help command prints help


	matro server  -f [--file] <file path>
	It generates all server side code
	[-vv] for verbose


	matro Relay  f [--file] <file path>
	It generates all Relay side code
	[vv] for verbose


	version
	matro version 

	`
	type fields struct {
		CommandWriter *CommandWriter
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantToWrite string
	}{
		{
			name: "shoud print help text for all commands",
			fields: fields{
				CommandWriter: NewCommandWriter(new(bytes.Buffer)),
			},
			wantErr:     false,
			wantToWrite: expectedOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Help{
				CommandWriter: tt.fields.CommandWriter,
			}
			if err := h.Exec(); (err != nil) != tt.wantErr {
				t.Errorf("Help.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
			if out := tt.fields.CommandWriter.writer.(*bytes.Buffer).String(); cleanStr(out) != cleanStr(tt.wantToWrite) {
				t.Errorf("Help.Exec() expects = \n%v, \nprinted \n%v", cleanStr(tt.wantToWrite), cleanStr(out))
			}
		})
	}
}

func cleanStr(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, "\n", ""), "\t", ""), " ", "")
}
