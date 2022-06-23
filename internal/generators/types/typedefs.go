package types

import (
	"fmt"

	"github.com/firstcontributions/matro/internal/errors"
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

// GetTypeDefs get all typedefs from the parsed json schema
func GetTypeDefs(d *parser.Definition) (*TypeDefs, error) {
	if d.Modules == nil {
		return nil, errors.ErrNoModules
	}
	types := []*CompositeType{}
	edges := utils.NewSet[string]()
	allTypesMap := getParsedTypesMap(d)
	queriesModule := parser.Module{
		Name: "queries",
	}
	queries, queryTypes, err := getQueries(d, allTypesMap, queriesModule)
	if err != nil {
		return nil, err
	}
	for _, m := range d.Modules {
		for _, def := range m.Entities {
			t, err := NewCompositeType(d, m, allTypesMap, def)
			if err != nil {
				return nil, err
			}
			edges.Union(t.EdgeFields())
			types = append(types, t)
			queries = append(queries, t.Queries()...)
		}
	}

	for _, q := range queries {
		edges.Add(q.Type)
	}
	typesMap := getTypeMap(d, types, edges)
	if _, ok := typesMap[d.Defaults.ViewerType]; !ok {
		return nil, fmt.Errorf("could not find viewer type [%s] in type definitions", d.Defaults.ViewerType)
	}
	return &TypeDefs{
		Types:      typesMap,
		Queries:    queries,
		QueryTypes: queryTypes,
	}, nil
}

// getTypeMap generated the <typeName><Type> map
func getTypeMap(d *parser.Definition, types []*CompositeType, edges *utils.Set[string]) map[string]*CompositeType {
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

// GetTypes gets list of types by name
func (g *TypeDefs) GetTypes(types map[string]*parser.Type) []*CompositeType {
	typeDefs := []*CompositeType{}
	for t := range types {
		typeDefs = append(typeDefs, g.Types[t])
	}
	return typeDefs
}
