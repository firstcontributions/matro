package commands

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/generators/relayjs"
	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/parser"
)

func TestRelay_getGenerators(t *testing.T) {
	d, _ := parser.NewDefinition().ParseFrom(strings.NewReader(validConfig))
	typeDefs, _ := types.GetTypeDefs(d)

	r := NewRelay(NewCommandWriter(new(bytes.Buffer)))

	if g := r.getGenerators(d, typeDefs); !reflect.DeepEqual(g, []generators.IGenerator{
		relayjs.NewGenerator(r.outputPath, d, typeDefs),
	}) {
		t.Errorf("s.getGenerators not matching expectation")
	}
}

func TestRelay_Exec(t *testing.T) {
	tests := []struct {
		name          string
		CodeGenerator *CodeGenerator
		wantErr       bool
	}{
		{
			name: "should return nil if help flaf is provided",
			CodeGenerator: &CodeGenerator{
				help:          true,
				CommandWriter: NewCommandWriter(new(bytes.Buffer)),
			},
			wantErr: false,
		},
		{
			name: "should throw error if parsing config failed",
			CodeGenerator: &CodeGenerator{
				CommandWriter: NewCommandWriter(new(bytes.Buffer)),
				filepath:      "random.txt",
			},
			wantErr: true,
		},
		{
			name: "should throw error if generators failed",
			CodeGenerator: &CodeGenerator{
				CommandWriter: NewCommandWriter(new(bytes.Buffer)),
				filepath:      "valid_config.json",
				outputPath:    "/",
			},
			wantErr: true,
		},
		{
			name: "should return nil if no errors",
			CodeGenerator: &CodeGenerator{
				CommandWriter: NewCommandWriter(new(bytes.Buffer)),
				filepath:      "valid_config.json",
				outputPath:    "/tmp/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Relay{
				CodeGenerator: tt.CodeGenerator,
			}
			c.CodeGenerator.FS = fsys

			if err := c.Exec(); (err != nil) != tt.wantErr {
				t.Errorf("Server.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
