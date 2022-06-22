package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCmd_Exec(t *testing.T) {
	versionText := "Version :  \nMinVersion :  \nBuildTime : \n"
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
			name: "should print the version details",
			fields: fields{
				CommandWriter: NewCommandWriter(new(bytes.Buffer)),
			},
			wantErr:     false,
			wantToWrite: versionText,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VersionCmd{
				CommandWriter: tt.fields.CommandWriter,
			}
			if err := v.Exec(); (err != nil) != tt.wantErr {
				t.Errorf("VersionCmd.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
			if out := tt.fields.CommandWriter.writer.(*bytes.Buffer).String(); strings.Replace(out, " ", "", -1) != strings.Replace(tt.wantToWrite, " ", "", -1) {
				t.Errorf("expected %v got %v", tt.wantToWrite, out)
			}
		})
	}
}
