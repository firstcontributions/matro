package commands

import (
	"bytes"
	"testing"
	"testing/fstest"

	"github.com/sirupsen/logrus"
)

var (
	invalidJson   = `}{`
	invalidConfig = `{}`
	validConfig   = `{
		"repo": "github.com/firstcontributions/test-matro",
		"defaults": {
			"viewer_type": "user"
		},
		"high_level_queries": [
		],
		"modules": [
			{
				"name": "users",
				"data_source": "mongo",
				"db": "mongo",
				"entities": {
					"user":{
						"name": "user",
						"type": "object",
						"meta": {
							"search_fields": ["name"],
							"filters": ["name"],
							"graphql_ops": "CRUD"
						},
						"properties":{
							"id": "id",
							"name": "string"
						}
					}
				}
			}
		]
	}`
)

func TestCodeGenerator_Setup_Help(t *testing.T) {
	helpText := `
	matro <code generator>  -f [--file] <file path>
	[-vv] for verbose

	`

	c := &CodeGenerator{
		CommandWriter: NewCommandWriter(new(bytes.Buffer)),
		help:          true,
	}
	if countinue := c.Setup(); countinue {
		t.Errorf("Server.Setup() continue = %v", countinue)
	}

	if out := c.CommandWriter.writer.(*bytes.Buffer).String(); out != helpText {
		t.Errorf("did not print help text %v %v", out, helpText)
	}
}

func TestCodeGenerator_Setup_LogLevel(t *testing.T) {

	tests := []struct {
		name         string
		verbose      bool
		wantLogLevel logrus.Level
		wantErr      bool
	}{
		{
			name:         "log level should be debug if verbose enabled",
			verbose:      true,
			wantLogLevel: logrus.DebugLevel,
		},
		{
			name:         "log level should be fatal if verbose enabled",
			verbose:      false,
			wantLogLevel: logrus.FatalLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CodeGenerator{
				verbose: tt.verbose,
			}
			c.Setup()
			if level := logrus.GetLevel(); level != tt.wantLogLevel {
				t.Errorf("Server.Exec() expected log level = %v, got %v", tt.wantLogLevel, level)
			}
		})
	}
}

func TestCodeGenerator_Exec(t *testing.T) {
	fs := fstest.MapFS{
		"valid_config.json": {
			Data: []byte(validConfig),
		},
		"invalid_config.json": {
			Data: []byte(invalidConfig),
		},
		"random.txt": {
			Data: []byte(invalidJson),
		},
	}

	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{
			name:    "should throw error if no config file exists",
			wantErr: true,
		},
		{
			name:     "should throw error if there is any error in parsing config",
			wantErr:  true,
			filepath: "random.txt",
		},
		{
			name:     "should throw error if there is any error in getting type defenitions",
			wantErr:  true,
			filepath: "invalid_config.json",
		},
		{
			name:     "should parse a valid input file",
			filepath: "valid_config.json",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CodeGenerator{
				filepath: tt.filepath,
				FS:       fs,
			}
			if _, _, err := c.GetDefenitionsAndTypes(); (err != nil) != tt.wantErr {
				t.Errorf("Server.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
