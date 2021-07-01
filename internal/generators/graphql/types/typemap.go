package types

import (
	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

var graphQLSchemaTypeMap = map[string]string{
	parser.ID:     "ID",
	parser.String: "String",
	parser.Time:   "Time",
	parser.Int:    "Int",
	parser.Float:  "Float",
	parser.Bool:   "Boolean",
}
var goTypeMap = map[string]string{
	parser.ID:     "graphql.ID",
	parser.String: "string",
	parser.Time:   "graphql.Time",
	parser.Int:    "int32",
	parser.Float:  "float64",
	parser.Bool:   "bool",
}

func getGraphQLType(t string) string {
	if s, ok := graphQLSchemaTypeMap[t]; ok {
		return s
	}
	return utils.ToTitleCase(t)
}
func getGoType(t string) string {
	if s, ok := goTypeMap[t]; ok {
		return s
	}
	return utils.ToTitleCase(t)
}
