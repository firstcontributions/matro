package types

import (
	"github.com/firstcontributions/matro/internal/parser"
)

// Query encapsulates a graphql query metadata
type Query struct {
	*Field
	Parent *CompositeType
}

func getQueries(d *parser.Definition, typesMap map[string]*parser.Type, queryModule parser.Module) ([]Query, map[string]*CompositeType) {
	queries := []Query{}
	graphQLOnlyTypes := map[string]*CompositeType{}
	for _, q := range d.Queries {
		if q.Schema == "" {
			t := NewCompositeType(typesMap, q, queryModule)
			graphQLOnlyTypes[t.Name] = t
			queries = append(queries, t.Queries()...)
		}
		field := NewField(typesMap, q, q.Name)
		field.IsQuery = true
		queries = append(queries, Query{
			Field: field,
		})
	}
	return queries, graphQLOnlyTypes
}
