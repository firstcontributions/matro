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
		"camel":    utils.ToCamelCase,
		"type":     types.GetGraphQLType,
		"grpcType": types.GetGRPCType,
		"add":      func(a, b int) int { return a + b },
		"counter":  utils.Counter(),
		"plural":   p.Plural,
		"empty": func(a []string) bool {
			return len(a) == 0
		},
		"isElemOfStrArray": utils.IsElementOfStringArray,
		"getTypeFromMap": func(m map[string]*types.CompositeType, t string) *types.CompositeType {
			return m[t]
		},
		"isCompositeType": types.IsCompositeType,
		"isHardCodedFilter": func(hardcodedFilters map[string]string, filter string) bool {
			_, ok := hardcodedFilters[filter]
			return ok
		},
		"getHardcodedValue": func(hardcodedFilters map[string]string, filter, typ string) string {
			val := hardcodedFilters[filter]
			if types.IsCompositeType(typ) || typ == "string" {
				return "\"" + val + "\""
			}
			return val
		},
		"isAditField": func(field string) bool {
			return field == "time_created" || field == "time_updated" || field == "id"
		},
	}
}
