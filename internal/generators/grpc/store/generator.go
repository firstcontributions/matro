package store

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
		if err := g.generateGRPCStore(ctx, m); err != nil {
			return err
		}
		for _, t := range m.Types {
			if t.IsNode {
				if err := g.generateModel(ctx, m.Name, t); err != nil {
					return err
				}
			}

		}
	}
	return nil
}

// generateMongoStore generates a mongo implementation for the store interface,
// constants associated, connection pool manager etc
func (g *Generator) generateGRPCStore(ctx context.Context, m Module) error {
	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/models/%sstore/grpc", g.Path, m.Name),
		"store.go",
		storeTmpl,
		m,
	)
}

// generateModel generates crud operations for the given types
// supported operations:
// Create, GetAll(search, filter, pagination),
// GetByID, Update, Delete
func (g *Generator) generateModel(ctx context.Context, module string, typ *types.CompositeType) error {

	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/models/%sstore/grpc", g.Path, module),
		typ.Name+".go",
		crudTmpl,
		struct {
			Module string
			*types.CompositeType
			Repo string
		}{
			Module:        module,
			CompositeType: typ,
			Repo:          g.Repo,
		},
	)

}
