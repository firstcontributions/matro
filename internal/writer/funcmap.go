package writer

import (
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/gertd/go-pluralize"
)

// FuncMap returns a map of functions to be used for template generation
func FuncMap() template.FuncMap {
	p := pluralize.NewClient()

	return template.FuncMap{
		"title":    utils.ToTitleCase,
		"type":     types.GetGraphQLType,
		"grpcType": types.GetGRPCType,
		"add":      func(a, b int) int { return a + b },
		"counter":  utils.Counter(),
		"plural":   p.Plural,
	}
}
