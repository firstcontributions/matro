package proto

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/firstcontributions/matro/internal/writer"
	"github.com/sirupsen/logrus"
)

// Generator implements gRPC protobuf generator
type Generator struct {
	*types.TypeDefs
	modules map[string]Module
	Path    string
}

// Module encapsulates module meta data and types in module
type Module struct {
	parser.Module
	Types []*types.CompositeType
}

// NewGenerator returns an instance of gRPC code generator
func NewGenerator(path string, d *parser.Definition, td *types.TypeDefs) *Generator {
	mods := map[string]Module{}
	for _, m := range d.Modules {
		if m.DataSource == "grpc" {
			mods[m.Name] = Module{
				Module: m,
				Types:  td.GetTypes(m.Entities),
			}
		}
	}
	return &Generator{
		TypeDefs: td,
		modules:  mods,
		Path:     path,
	}
}

// Generate generates gRPC prtobuf code for all given modules(services)
func (g *Generator) Generate(ctx context.Context) error {
	for _, m := range g.modules {
		if m.DataSource != "grpc" {
			continue
		}
		if err := g.generateProtoForModule(ctx, m, tmpl); err != nil {
			return err
		}
		if err := g.generateGRPCService(g.protoFilePathForModule(m)); err != nil {
			return err
		}
	}
	return nil
}

// protoFilePathForModule returns path to protofile to be generated
func (g *Generator) protoFilePathForModule(m Module) string {
	return fmt.Sprintf("%s/api/%s.proto", g.Path, m.Name)
}

// generateProtoForModule generates proto file for the given module
func (g *Generator) generateProtoForModule(ctx context.Context, m Module, tmpl string) error {
	path := fmt.Sprintf("%s/api", g.Path)
	return writer.CompileAndWrite(
		ctx,
		path,
		m.Name+".proto",
		tmpl,
		m,
	)
}

// generateGRPCService generates grpc service stub from the proto file
func (*Generator) generateGRPCService(protoPath string) error {
	logrus.Info("going to generate grpc")
	if res, err := exec.Command(
		"protoc",
		"--go_out=.",
		"--go-grpc_out=.",
		protoPath,
	).Output(); err != nil {
		logrus.Debug("grpc error", string(res), err)
		return err
	}
	return nil
}
