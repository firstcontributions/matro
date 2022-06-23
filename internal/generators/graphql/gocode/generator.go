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
func NewGenerator(path string, d *parser.Definition, td *types.TypeDefs) *Generator {
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
		if t.NoGraphql {
			continue
		}
		if err := g.generateTypeResolver(ctx, t); err != nil {
			return err
		}
		if err := g.generateMutation(ctx, t); err != nil {
			return err
		}
	}
	for _, t := range g.QueryTypes {
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
			Repo  string
			Types map[string]*types.CompositeType
		}{
			CompositeType: t,
			Repo:          g.Repo,
			Types:         g.Types,
		},
	)
}

// generateMutation
func (g *Generator) generateMutation(ctx context.Context, t *types.CompositeType) error {
	if !t.IsNode || !(t.GraphqlOps.Create() || t.GraphqlOps.Update() || t.GraphqlOps.Delete()) {
		return nil
	}
	return writer.CompileAndWrite(
		ctx,
		g.Path,
		t.Name+"_mutations.go",
		mutationTmpl,
		t,
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
	var returnType *types.CompositeType

	typeName := q.Type
	if q.Type == "object" {
		typeName = q.Name
	}
	if t, ok := g.Types[typeName]; ok {
		returnType = t
	} else if t, ok := g.QueryTypes[typeName]; ok {
		returnType = t
	} else {
		return fmt.Errorf("could not find type defenition for type [%s]", typeName)
	}
	filename := q.Name + "queryresolver.go"
	if q.Parent != nil {
		filename = q.Parent.Name + filename
	}
	return writer.CompileAndWrite(
		ctx,
		g.Path,
		filename,
		queryResolverTmpl,
		struct {
			Query      types.Query
			ReturnType *types.CompositeType
		}{
			Query:      q,
			ReturnType: returnType,
		},
	)
}
