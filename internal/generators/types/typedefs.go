package types

import (
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/gertd/go-pluralize"
)

type TypeDefs struct {
	Types map[string]*CompositeType
}

func NewTypeDefs(path string, d *parser.Definition) *TypeDefs {
	types := []*CompositeType{}
	edges := utils.NewSet()
	for _, def := range d.DataSchema {
		t := NewCompositeType(d, def)
		edges.Union(t.EdgeTypes())
		types = append(types, t)
	}
	return &TypeDefs{
		Types: getTypeMap(types, edges),
	}
}
func getTypeMap(types []*CompositeType, edges *utils.Set) map[string]*CompositeType {
	typeMap := map[string]*CompositeType{}
	for _, t := range types {
		if edges.IsElem(t.Name) {
			t.IsEdge = true
		}
		typeMap[t.Name] = t
	}
	for _, t := range types {
		for _, f := range t.Fields {
			if f.IsJoinedData && f.IsList {
				typeMap[f.Type].ReferedFields = append(
					typeMap[f.Type].ReferedFields,
					t.Name,
				)

			}
		}
	}
	return typeMap
}

func (g *TypeDefs) FuncMap() template.FuncMap {
	p := pluralize.NewClient()

	return template.FuncMap{
		"title":    utils.ToTitleCase,
		"type":     getGraphQLType,
		"grpcType": getGRPCType,
		"add":      func(a, b int) int { return a + b },
		"counter":  utils.Counter(),
		"plural":   p.Plural,
	}
}

func (g *TypeDefs) GetTypeDefs(strTypes []string) []*CompositeType {
	typeDefs := []*CompositeType{}
	for _, t := range strTypes {
		typeDefs = append(typeDefs, g.Types[t])
	}
	return typeDefs
}
