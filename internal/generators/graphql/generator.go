package graphql

import (
	"github.com/firstcontributions/matro/internal/generators/graphql/gocode"
	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/generators/graphql/types"
	"github.com/firstcontributions/matro/internal/parser"
)

type IGenerator interface {
	Generate() error
}

func GetGenerator(path, t string, d *parser.Definition) IGenerator {
	td := types.NewTypeDefs(path, d)
	switch t {
	case "schema":
		return &schema.Generator{
			TypeDefs: td,
		}
	case "gocode":
		return &gocode.Generator{
			TypeDefs: td,
		}
	}
	return nil
}
