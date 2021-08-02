package gocode

import (
	"context"
	"fmt"

	"github.com/firstcontributions/matro/internal/parser"
	"github.com/firstcontributions/matro/internal/writer"

	"github.com/firstcontributions/matro/internal/generators/types"
)

// Generator implements graphql server code generator
type Generator struct {
	*types.TypeDefs
	Path    string
	Modules []parser.Module
	Repo    string
}

// NewGenerator returns an instance of graphql server code generator
func NewGenerator(path string, d *parser.Definition) *Generator {
	td := types.NewTypeDefs(path, d)
	return &Generator{
		Path:     path,
		TypeDefs: td,
		Repo:     d.Repo,
		Modules:  d.Modules,
	}
}

// Generate generates all graphql server codes
// (type definitons, query resolvers, mutation executions, ...)
func (g *Generator) Generate(ctx context.Context) error {
	path := fmt.Sprintf("%s/internal/graphql/schema", g.Path)
	for _, t := range g.Types {
		if err := g.generateTypes(ctx, typesTpl, t, path, t.Name+".go"); err != nil {
			return err
		}
	}
	if err := g.generateNodeResolver(ctx); err != nil {
		return nil
	}
	return g.generateRootResolver(ctx)
}

// generateTypes generate types based on the given template
func (g *Generator) generateTypes(ctx context.Context, tmpl string, data *types.CompositeType, path, filename string) error {
	return writer.CompileAndWrite(
		ctx,
		path,
		filename,
		tmpl,
		struct {
			*types.CompositeType
			Repo string
		}{
			CompositeType: data,
			Repo:          g.Repo,
		},
	)
}

// generateTypes generate types based on the given template
func (g *Generator) generateNodeResolver(ctx context.Context) error {
	path := fmt.Sprintf("%s/internal/graphql/schema", g.Path)
	return writer.CompileAndWrite(
		ctx,
		path,
		"node.go",
		nodeTmpl,
		g,
	)
}

// generateRootResolver generate root query resolver
func (g *Generator) generateRootResolver(ctx context.Context) error {
	path := fmt.Sprintf("%s/internal/graphql/schema", g.Path)
	return writer.CompileAndWrite(
		ctx,
		path,
		"resolver.go",
		resolverTmpl,
		g,
	)
}
