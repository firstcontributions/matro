package types

import (
	"fmt"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

// Query encapsulates a graphql query metadata
type Query struct {
	*Field
	Parent *CompositeType
}

func (q Query) InputName() string {
	if q.Parent == nil {
		return fmt.Sprintf("%sInput", utils.ToTitleCase(q.Name))
	}
	return fmt.Sprintf("%s%sInput", utils.ToTitleCase(q.Parent.Name), utils.ToTitleCase(q.Name))
}

func getQueries(d *parser.Definition, typesMap map[string]*parser.Type, queryModule parser.Module) ([]Query, map[string]*CompositeType, error) {
	queries := []Query{}
	graphQLOnlyTypes := map[string]*CompositeType{}
	for _, q := range d.Queries {
		if q.Schema == "" {
			t, err := NewCompositeType(d, queryModule, typesMap, q)
			if err != nil {
				return nil, nil, err
			}
			graphQLOnlyTypes[t.Name] = t
			queries = append(queries, t.Queries()...)
		}
		field := NewField(d, typesMap, q, q.Name, false)
		field.IsQuery = true
		queries = append(queries, Query{
			Field: field,
		})
	}
	return queries, graphQLOnlyTypes, nil
}
