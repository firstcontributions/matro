package schema

import (
	"context"
	"fmt"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/firstcontributions/matro/internal/writer"
)

// Generator implements the graphql schema generator
type Generator struct {
	*types.TypeDefs
	Path string
}

// NewGenerator returns an instance of graphql schema generator
func NewGenerator(path string, d *parser.Definition, td *types.TypeDefs) *Generator {
	return &Generator{
		Path:     path,
		TypeDefs: td,
	}
}

// Generate generates a graphql schema file based on given template
func (g *Generator) Generate(ctx context.Context) error {
	path := fmt.Sprintf("%s/assets", g.Path)
	return writer.CompileAndWrite(
		ctx,
		path,
		"schema.graphql",
		schemaTmpl,
		g,
	)
}
