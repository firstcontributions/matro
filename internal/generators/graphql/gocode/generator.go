package gocode

import (
	"fmt"

	"github.com/firstcontributions/matro/internal/parser"
	"github.com/firstcontributions/matro/internal/writer"

	"github.com/firstcontributions/matro/internal/generators/types"
)

// Generator implements graphql server code generator
type Generator struct {
	*types.TypeDefs
	Path string
}

// NewGenerator returns an instance of graphql server code generator
func NewGenerator(path string, d *parser.Definition) *Generator {
	td := types.NewTypeDefs(path, d)
	return &Generator{
		Path:     path,
		TypeDefs: td,
	}
}

// Generate generates all graphql server codes
// (type definitons, query resolvers, mutation executions, ...)
func (g *Generator) Generate() error {
	path := fmt.Sprintf("%s/internal/graphql/schema", g.Path)
	for _, t := range g.Types {
		if err := g.generateTypes(typesTpl, t, path, t.Name+".go"); err != nil {
			return err
		}
	}
	return nil
}

// generateTypes generate types based on the given template
func (g *Generator) generateTypes(tmpl string, data interface{}, path, filename string) error {
	return writer.CompileAndWrite(
		path,
		filename,
		tmpl,
		data,
	)
}
