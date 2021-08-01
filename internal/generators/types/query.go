package types

import "github.com/firstcontributions/matro/internal/parser"

// Query encapsulates a graphql query metadata
type Query struct {
	*Field
}

func getQueries(d *parser.Definition) []Query {
	queries := []Query{}
	for _, q := range d.Queries {
		field := NewField(d, q, q.Name)
		field.IsQuery = true
		queries = append(queries, Query{
			Field: field,
		})
	}
	return queries
}
