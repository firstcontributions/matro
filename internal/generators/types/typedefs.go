package types

import (
	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

// TypeDefs encapsulates the list of types
// it keeps the information in the form of a map
// of type name to Type struct
type TypeDefs struct {
	Types   map[string]*CompositeType
	Queries []Query
}

// NewTypeDefs get all typedefs from the parsed json schema
func NewTypeDefs(path string, d *parser.Definition) *TypeDefs {
	types := []*CompositeType{}
	edges := utils.NewSet()
	for _, def := range d.DataSchema {
		t := NewCompositeType(d, def)
		edges.Union(t.EdgeFields())
		types = append(types, t)
	}
	queries := getQueries(d)
	for _, q := range queries {
		edges.Add(q.Type)
	}
	return &TypeDefs{
		Types:   getTypeMap(d, types, edges),
		Queries: queries,
	}
}

// getTypeMap generated the <typeName><Type> map
func getTypeMap(d *parser.Definition, types []*CompositeType, edges *utils.Set) map[string]*CompositeType {
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
	for _, m := range d.Modules {
		for _, t := range m.Entities {
			if _, ok := typeMap[t]; ok {
				typeMap[t].Module = m.Name
			}
		}
	}
	return typeMap
}

// GetTypeDefs gets list of types by name
func (g *TypeDefs) GetTypeDefs(strTypes []string) []*CompositeType {
	typeDefs := []*CompositeType{}
	for _, t := range strTypes {
		typeDefs = append(typeDefs, g.Types[t])
	}
	return typeDefs
}
