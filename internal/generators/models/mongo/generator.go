package mongo

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
	Modules map[string]Module
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
		Modules:  mods,
		Path:     path,
		Repo:     d.Repo,
	}
}

// Generate generates a store interface, mongo implementation for the interface,
// data schema, crud operations for all the types in all given modules
func (g *Generator) Generate(ctx context.Context) error {
	if err := g.generateCursorPkg(ctx); err != nil {
		return err
	}
	if err := g.generateUtilsPkg(ctx); err != nil {
		return err
	}
	if err := g.generateDataProcessorInterface(ctx); err != nil {
		return err
	}
	for _, m := range g.Modules {
		if err := g.generateStore(ctx, m); err != nil {
			return err
		}
		for _, t := range m.Types {
			if t.AllReferedFields {
				continue
			}
			if err := g.generateModelTypes(ctx, m, t); err != nil {
				return err
			}
		}
		if m.DB == "" {
			// if there is no DB information, not needed to generate mongo related
			continue
		}
		if err := g.generateMongoStore(ctx, m); err != nil {
			return err
		}
		for _, t := range m.Types {
			if t.AllReferedFields {
				continue
			}
			if t.IsNode {
				if err := g.generateModel(ctx, t); err != nil {
					return err
				}
			}
		}
	}
	return g.generateStoreManage(ctx)
}

func (g *Generator) generateStoreManage(ctx context.Context) error {
	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/storemanager/", g.Path),
		"storemanager.go",
		storeManagerTmpl,
		g,
	)
}

func (g *Generator) generateCursorPkg(ctx context.Context) error {
	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/pkg/cursor/", g.Path),
		"cursor.go",
		cursorTmpl,
		g,
	)
}

func (g *Generator) generateUtilsPkg(ctx context.Context) error {
	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/models/utils", g.Path),
		"pagination.go",
		utilsTmpl,
		g,
	)
}

func (g *Generator) generateDataProcessorInterface(ctx context.Context) error {
	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/models/utils", g.Path),
		"dataprocessor.go",
		dprocTmpl,
		g,
	)
}

// generateStore generates  the store interface,
func (g *Generator) generateStore(ctx context.Context, m Module) error {
	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/models/%sstore", g.Path, m.Name),
		"store.go",
		storeInterfaceTpl,
		m,
	)
}

// generateMongoStore generates a mongo implementation for the store interface,
// constants associated, connection pool manager etc
func (g *Generator) generateMongoStore(ctx context.Context, m Module) error {
	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/models/%sstore/mongo", g.Path, m.Name),
		"store.go",
		storeTpl,
		m,
	)
}

// generateModel generates crud operations for the given types
// supported operations:
// Create, GetAll(search, filter, pagination),
// GetByID, Update, Delete
func (g *Generator) generateModel(ctx context.Context, typ *types.CompositeType) error {

	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/models/%sstore/mongo", g.Path, typ.Module.Name),
		typ.Name+".go",
		crudTpl,
		struct {
			*types.CompositeType
			Repo string
		}{
			CompositeType: typ,
			Repo:          g.Repo,
		},
	)

}

// generateModelTypes generates data schema for the given types
func (g *Generator) generateModelTypes(ctx context.Context, module Module, typ *types.CompositeType) error {
	return writer.CompileAndWrite(
		ctx,
		fmt.Sprintf("%s/internal/models/%sstore", g.Path, module.Name),
		typ.Name+".go",
		modelTyp,
		struct {
			Module Module
			*types.CompositeType
			HaveProto bool
		}{
			Module:        module,
			CompositeType: typ,
			HaveProto:     module.DataSource == "grpc",
		},
	)
}
