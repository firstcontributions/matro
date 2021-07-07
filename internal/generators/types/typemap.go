package types

import (
	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

// graphQLSchemaTypeMap is a mapping between matro types
// and graphql types
var graphQLSchemaTypeMap = map[string]string{
	parser.ID:     "ID",
	parser.String: "String",
	parser.Time:   "Time",
	parser.Int:    "Int",
	parser.Float:  "Float",
	parser.Bool:   "Boolean",
}

// graphQLSchemaTypeMap is a mapping between matro types and
// graphql go implementation types

var goGraphQLTypeMap = map[string]string{
	parser.ID:     "graphql.ID",
	parser.String: "string",
	parser.Time:   "graphql.Time",
	parser.Int:    "int32",
	parser.Float:  "float64",
	parser.Bool:   "bool",
}

// graphQLSchemaTypeMap is a mapping between matro types and
// golang types
var goTypeMap = map[string]string{
	parser.ID:     "string",
	parser.String: "string",
	parser.Time:   "time.Time",
	parser.Int:    "int32",
	parser.Float:  "float64",
	parser.Bool:   "bool",
}

// graphQLSchemaTypeMap is a mapping between matro types and
// protobuf types
var grpcTypeMap = map[string]string{
	parser.ID:     "string",
	parser.String: "string",
	parser.Time:   "google.protobuf.Timestamp",
	parser.Int:    "int32",
	parser.Float:  "float64",
	parser.Bool:   "bool",
}

// getGraphQLType returns graphql schema type from matro type
func getGraphQLType(t string) string {
	if s, ok := graphQLSchemaTypeMap[t]; ok {
		return s
	}
	return utils.ToTitleCase(t)
}

// getGoType returns go type from matro type
func getGoType(t string) string {
	if s, ok := goTypeMap[t]; ok {
		return s
	}
	return utils.ToTitleCase(t)
}

// getGoGraphQLType returns graphql go implementation
//  type from matro type
func getGoGraphQLType(t string) string {
	if s, ok := goGraphQLTypeMap[t]; ok {
		return s
	}
	return utils.ToTitleCase(t)
}

// getGRPCType protobuf type from matro type
func getGRPCType(t string) string {
	if s, ok := grpcTypeMap[t]; ok {
		return s
	}
	return utils.ToTitleCase(t)
}
