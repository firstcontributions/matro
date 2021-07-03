package proto

import (
	"fmt"
	"os/exec"
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

type Generator struct {
	*types.TypeDefs
	modules map[string]Module
	Path    string
}

type Module struct {
	parser.Module
	Types []*types.CompositeType
}

func NewGenerator(path string, d *parser.Definition) *Generator {
	td := types.NewTypeDefs(path, d)
	mods := map[string]Module{}
	for _, m := range d.Modules {
		if m.DataSource == "grpc" {
			mods[m.Name] = Module{
				Module: m,
				Types:  td.GetTypeDefs(m.Entities),
			}
		}
	}
	return &Generator{
		TypeDefs: td,
		modules:  mods,
		Path:     path,
	}
}
func (g *Generator) Generate() error {
	path := fmt.Sprintf("%s/api", g.Path)
	t, err := template.New("proto").
		Funcs(g.FuncMap()).
		Parse(tmpl)
	if err != nil {
		return err
	}

	for name, m := range g.modules {
		fw, err := utils.GetFileWriter(path, name+".proto")
		if err != nil {
			return err
		}
		if err := t.Execute(fw, m); err != nil {
			return err
		}
		if err := fw.Close(); err != nil {
			return err
		}
		if err := g.generateGRPCService(fmt.Sprintf("%s/%s.proto", path, name)); err != nil {
			return nil
		}
	}
	return nil
}

func (g *Generator) generateGRPCService(protoPath string) error {
	if _, err := exec.Command("protoc", protoPath, "--go_out=./__generated", "--go-grpc_out=./__generated").Output(); err != nil {
		return err
	}
	return nil
}
