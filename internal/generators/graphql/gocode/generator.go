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
		Path:     fmt.Sprintf("%s/internal/graphql/schema", path),
		TypeDefs: td,
		Repo:     d.Repo,
		Modules:  d.Modules,
	}
}

// Generate generates all graphql server codes
// (type definitons, query resolvers, mutation executions, ...)
func (g *Generator) Generate(ctx context.Context) error {
	for _, t := range g.Types {
		if err := g.generateTypeResolver(ctx, t); err != nil {
			return err
		}
	}
	for _, q := range g.Queries {
		if err := g.generateQueryResolver(ctx, q); err != nil {
			return err
		}
	}
	if err := g.generateNodeResolver(ctx); err != nil {
		return nil
	}
	return g.generateRootResolver(ctx)
}

// generateTypeResolver generate types based on the given template
func (g *Generator) generateTypeResolver(ctx context.Context, t *types.CompositeType) error {
	return writer.CompileAndWrite(
		ctx,
		g.Path,
		t.Name+"resolver.go",
		typesTpl,
		struct {
			*types.CompositeType
			Repo string
		}{
			CompositeType: t,
			Repo:          g.Repo,
		},
	)
}

// generateTypes generate types based on the given template
func (g *Generator) generateNodeResolver(ctx context.Context) error {
	return writer.CompileAndWrite(
		ctx,
		g.Path,
		"noderesolver.go",
		nodeTmpl,
		g,
	)
}

// generateRootResolver generate root query resolver
func (g *Generator) generateRootResolver(ctx context.Context) error {
	return writer.CompileAndWrite(
		ctx,
		g.Path,
		"resolver.go",
		resolverTmpl,
		g,
	)
}

// generateQueryResolver generate root query resolver
func (g *Generator) generateQueryResolver(ctx context.Context, q types.Query) error {
	return writer.CompileAndWrite(
		ctx,
		g.Path,
		q.Name+"resolver.go",
		queryResolverTmpl,
		struct {
			Query      types.Query
			ReturnType *types.CompositeType
		}{
			Query:      q,
			ReturnType: g.Types[q.Type],
		},
	)
}
