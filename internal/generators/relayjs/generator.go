package relayjs

import (
	"context"
	"fmt"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/firstcontributions/matro/internal/writer"
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
		Path:     fmt.Sprintf("%s/src/components/", path),
		TypeDefs: td,
		Repo:     d.Repo,
		Modules:  d.Modules,
	}
}

func (g *Generator) Generate(ctx context.Context) error {
	for _, q := range g.Queries {
		if q.Parent == nil || !q.IsPaginated {
			continue
		}
		if err := g.generatePaginatedFragments(ctx, q); err != nil {
			return err
		}
	}
	for _, t := range g.Types {
		if err := g.generateTypes(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) generatePaginatedFragments(ctx context.Context, query types.Query) error {

	return writer.CompileAndWrite(
		ctx,
		g.Path+query.Type,
		fmt.Sprintf("%sList.tsx", utils.ToTitleCase(query.Type)),
		paginationFragment,
		query,
	)
}

func (g *Generator) generateTypes(ctx context.Context, t *types.CompositeType) error {

	return writer.CompileAndWrite(
		ctx,
		g.Path+t.Name,
		fmt.Sprintf("%s.tsx", utils.ToTitleCase(t.Name)),
		typesTmpl,
		t,
	)
}
