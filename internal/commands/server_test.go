package commands

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/generators/gomod"
	"github.com/firstcontributions/matro/internal/generators/graphql/gocode"
	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/generators/grpc/proto"
	"github.com/firstcontributions/matro/internal/generators/grpc/service"
	"github.com/firstcontributions/matro/internal/generators/grpc/store"
	"github.com/firstcontributions/matro/internal/generators/models/mongo"
	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/parser"
)

func TestServer_getGenerators(t *testing.T) {
	d, _ := parser.NewDefinition().ParseFrom(strings.NewReader(validConfig))
	typeDefs, _ := types.GetTypeDefs(d)

	s := NewServer(NewCommandWriter(new(bytes.Buffer)))

	if g := s.getGenerators(d, typeDefs); !reflect.DeepEqual(g, []generators.IGenerator{
		gomod.NewGenerator(s.outputPath, d),
		schema.NewGenerator(s.outputPath, d, typeDefs),
		gocode.NewGenerator(s.outputPath, d, typeDefs),
		proto.NewGenerator(s.outputPath, d, typeDefs),
		store.NewGenerator(s.outputPath, d, typeDefs),
		mongo.NewGenerator(s.outputPath, d, typeDefs),
		service.NewGenerator(s.outputPath, d, typeDefs),
	}) {
		t.Errorf("s.getGenerators not matching expectation")
	}

}

func TestServer_Exec(t *testing.T) {
	path := os.Getenv("MATRO_PATH")
	if path == "" {
		path = "/tmp/"
	}
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
				outputPath:    path,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Server{
				CodeGenerator: tt.CodeGenerator,
			}
			c.CodeGenerator.FS = fsys

			if err := c.Exec(); (err != nil) != tt.wantErr {
				t.Errorf("Server.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
