package mongo

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
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
	parser.Module
	Types []*types.CompositeType
}

// NewGenerator returns an instance of mongo model generator
func NewGenerator(path string, d *parser.Definition) *Generator {
	td := types.NewTypeDefs(path, d)
	mods := map[string]Module{}
	for _, m := range d.Modules {
		mods[m.Name] = Module{
			Module: m,
			Types:  td.GetTypeDefs(m.Entities),
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
func (g *Generator) Generate() error {
	for _, m := range g.modules {
		if err := g.generateStore(m); err != nil {
			return err
		}
		for _, t := range m.Types {
			if err := g.generateModel(m.Name, t); err != nil {
				return err
			}
			if err := g.generateModelTypes(m.Name, t); err != nil {
				return err
			}
		}
	}
	return nil
}

// generateStore generates a mongo implementation for the store interface,
// constants associated, connection pool manager etc
func (g *Generator) generateStore(m Module) error {
	t, err := template.New("mongo_store").
		Funcs(g.FuncMap()).
		Parse(storeTpl)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	if err := t.Execute(&b, m); err != nil {
		return err
	}
	return utils.FormatAndWriteGoCode(
		fmt.Sprintf("%s/internal/models/%sstore/mongo", g.Path, m.Name),
		"store.go",
		b.Bytes(),
	)
}

// generateModel generates crud operations for the given types
// supported operations:
// Create, GetAll(search, filter, pagination),
// GetByID, Update, Delete
func (g *Generator) generateModel(module string, typ *types.CompositeType) error {
	t, err := template.New("mongo_model").
		Funcs(g.FuncMap()).
		Parse(modelTpl)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	if err := t.Execute(
		&b,
		struct {
			Module string
			*types.CompositeType
			Repo string
		}{
			Module:        module,
			CompositeType: typ,
			Repo:          g.Repo,
		},
	); err != nil {
		return err
	}
	return utils.FormatAndWriteGoCode(
		fmt.Sprintf("%s/internal/models/%sstore/mongo", g.Path, module),
		typ.Name+".go",
		b.Bytes(),
	)
}

// generateModelTypes generates data schema for the given types
func (g *Generator) generateModelTypes(module string, typ *types.CompositeType) error {
	t, err := template.New("mongo_model_type").
		Funcs(g.FuncMap()).
		Parse(modelTyp)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	if err := t.Execute(
		&b,
		struct {
			Module string
			*types.CompositeType
		}{
			Module:        module,
			CompositeType: typ,
		},
	); err != nil {
		return err
	}
	return utils.FormatAndWriteGoCode(
		fmt.Sprintf("%s/internal/models/%sstore", g.Path, module),
		typ.Name+".go",
		b.Bytes(),
	)
}
