package service

import (
	"context"
	"fmt"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/firstcontributions/matro/internal/writer"
)

// Generator implements mongo model code generator
type Generator struct {
	*types.TypeDefs
	modules map[string]Module
	Path    string
	Repo    string
}

// Module encapsulates module metadata and types associated
type Module struct {
	Repo string
	parser.Module
	Types []*types.CompositeType
}

// NewGenerator returns an instance of mongo model generator
func NewGenerator(path string, d *parser.Definition, td *types.TypeDefs) *Generator {
	mods := map[string]Module{}
	for _, m := range d.Modules {
		mods[m.Name] = Module{
			Module: m,
			Types:  td.GetTypes(m.Entities),
			Repo:   d.Repo,
		}

	}
	return &Generator{
		TypeDefs: td,
		modules:  mods,
		Path:     path,
		Repo:     d.Repo,
	}
}

// Generate generates a store interface, mongo implementation for the interface,
// data schema, crud operations for all the types in all given modules
func (g *Generator) Generate(ctx context.Context) error {
	for _, m := range g.modules {
		if m.DataSource != "grpc" {
			continue
		}
		for _, t := range m.Types {
			if t.IsNode {
				if err := g.generateModel(ctx, t); err != nil {
					return err
				}
			}

		}
	}
	return nil
}

// generateModel generates crud operations for the given types
// supported operations:
// Create, GetAll(search, filter, pagination),
// GetByID, Update, Delete
func (g *Generator) generateModel(ctx context.Context, t *types.CompositeType) error {

	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/grpc/%s/service/", g.Path, t.Module.Name),
		t.Name+".go",
		crudTmpl,
		t,
	)

}
