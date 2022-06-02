package types

import (
	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

// TypeDefs encapsulates the list of types
// it keeps the information in the form of a map
// of type name to Type struct
type TypeDefs struct {
	Types      map[string]*CompositeType
	Queries    []Query
	QueryTypes map[string]*CompositeType
}

func getParsedTypesMap(d *parser.Definition) map[string]*parser.Type {
	parsedTypesMap := map[string]*parser.Type{}
	for _, m := range d.Modules {
		for name, def := range m.Entities {
			parsedTypesMap[name] = def
		}
	}
	return parsedTypesMap
}

// NewTypeDefs get all typedefs from the parsed json schema
func NewTypeDefs(path string, d *parser.Definition) *TypeDefs {
	types := []*CompositeType{}
	edges := utils.NewSet()
	allTypesMap := getParsedTypesMap(d)
	queriesModule := parser.Module{
		Name: "queries",
	}
	queries, queryTypes := getQueries(d, allTypesMap, queriesModule)
	for _, m := range d.Modules {
		for _, def := range m.Entities {
			t := NewCompositeType(allTypesMap, def, m)
			edges.Union(t.EdgeFields())
			types = append(types, t)
			queries = append(queries, t.Queries()...)
		}
	}

	for _, q := range queries {
		edges.Add(q.Type)
	}
	return &TypeDefs{
		Types:      getTypeMap(d, types, edges),
		Queries:    queries,
		QueryTypes: queryTypes,
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
			if f.IsJoinedData && f.IsList && !t.AllReferedFields {
				typeMap[f.Type].ReferedTypes[t.Name] = t
			}
			if _, ok := typeMap[f.Type]; ok {
				typeMap[f.Type].ParentTypes[t.Name] = t
			}
		}
	}
	for _, t := range typeMap {
		if t.NoGraphql || t.IsNode {
			continue
		}
		for _, rt := range t.ParentTypes {
			typeMap[t.Name].GraphqlOps.Union(rt.GraphqlOps)
		}
	}
	return typeMap
}

// GetTypeDefs gets list of types by name
func (g *TypeDefs) GetTypeDefs(types map[string]*parser.Type) []*CompositeType {
	typeDefs := []*CompositeType{}
	for t := range types {
		typeDefs = append(typeDefs, g.Types[t])
	}
	return typeDefs
}
