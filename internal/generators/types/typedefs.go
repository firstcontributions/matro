package types

import (
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

type TypeDefs struct {
	Types []*CompositeType
	Path  string
}

func NewTypeDefs(path string, d *parser.Definition) *TypeDefs {
	types := []*CompositeType{}
	edges := utils.NewSet()
	for _, def := range d.DataSchema {
		t := NewCompositeType(d, def)
		edges.Union(t.EdgeTypes())
		types = append(types, t)
	}
	for _, t := range types {
		if edges.IsElem(t.Name) {
			t.IsEdge = true
		}
	}
	return &TypeDefs{
		Types: types,
		Path:  path,
	}
}

func (g *TypeDefs) FuncMap() template.FuncMap {
	return template.FuncMap{
		"title": utils.ToTitleCase,
		"type":  getGraphQLType,
	}
}
