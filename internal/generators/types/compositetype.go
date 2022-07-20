package types

import (
	"fmt"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/gertd/go-pluralize"
)

var auditFields = []*Field{
	{
		Name:        "time_created",
		Type:        "time",
		IsPrimitive: true,
	},
	{
		Name:        "time_updated",
		Type:        "time",
		IsMutatable: true,
		IsPrimitive: true,
	},
}

// CompositeType defines a non trivial type with a set
// of informations like what all fields, opeations supported, etc.
type CompositeType struct {
	Name               string
	Fields             map[string]*Field
	ReferedTypes       map[string]*CompositeType
	ParentTypes        map[string]*CompositeType
	IsNode             bool
	IsEdge             bool
	Filters            []string
	GraphqlOps         parser.Ops
	SearchFields       []string
	MutatableFields    []string
	SortBy             []string
	Module             *parser.Module
	AllReferedFields   bool
	HardcodedFilters   map[string]string
	NoGraphql          bool
	IsViewerType       bool
	ViewerRefenceField string
}

// NewCompositeType return an instance of the CompositeType
func NewCompositeType(d *parser.Definition, module parser.Module, typesMap map[string]*parser.Type, typeDef *parser.Type) (*CompositeType, error) {
	fields := map[string]*Field{}
	isNode := false
	allRefered := true
	for field, def := range typeDef.Properties {
		if field == "id" {
			isNode = true
		}
		f := NewField(d, typesMap, def, field, typeDef.Meta.ViewerRefenceField == field)
		fields[field] = f
		if !(f.IsJoinedData && f.IsList) {
			allRefered = false
		}
	}
	if module.DB != "" && isNode {
		for _, f := range auditFields {
			fields[f.Name] = f
		}
	}
	for _, mf := range typeDef.Meta.MutatableFields {
		fields[mf].IsMutatable = true
	}
	sortBy := utils.NewSet(typeDef.Meta.SortBy...)
	if _, ok := fields["time_created"]; ok {
		sortBy.Add("time_created")
	}

	filters := utils.NewSet(typeDef.Meta.Filters...)
	if typeDef.Meta.ViewerRefenceField != "" {
		filters.Add(typeDef.Meta.ViewerRefenceField)
	}

	return &CompositeType{
		Name:               typeDef.Name,
		Fields:             fields,
		IsNode:             isNode,
		GraphqlOps:         typeDef.Meta.GraphqlOps,
		SearchFields:       typeDef.Meta.SearchFields,
		Filters:            filters.Elems(),
		MutatableFields:    typeDef.Meta.MutatableFields,
		SortBy:             sortBy.Elems(),
		Module:             &module,
		AllReferedFields:   allRefered,
		NoGraphql:          typeDef.NoGraphql,
		ReferedTypes:       map[string]*CompositeType{},
		ParentTypes:        map[string]*CompositeType{},
		IsViewerType:       typeDef.Name == d.Defaults.ViewerType,
		ViewerRefenceField: typeDef.Meta.ViewerRefenceField,
	}, nil
}

func (c *CompositeType) Queries() []Query {
	queries := []Query{}
	for _, f := range c.Fields {
		if f.IsQuery {
			queries = append(queries, Query{
				Field:  f,
				Parent: c,
			})
		}
	}
	return queries
}

// Mutatable will say this type is mutatable for not
func (c *CompositeType) Mutatable() bool {
	return len(c.MutatableFields) > 0
}

// EdgeFields return the paginated fields that can be an edge
func (c *CompositeType) EdgeFields() *utils.Set[string] {
	s := utils.NewSet[string]()
	for _, f := range c.Fields {
		if f.IsPaginated && f.IsList {
			s.Add(f.Type)
		}
	}
	return s
}

// EdgeName returns the edge type name
func (c *CompositeType) EdgeName() string {
	return fmt.Sprintf("%sEdge", utils.ToTitleCase(c.Name))
}

// ConnectionName returns the connection type name
func (c *CompositeType) ConnectionName() string {
	pl := pluralize.NewClient().Plural(c.Name)
	return fmt.Sprintf("%sConnection", utils.ToTitleCase(pl))
}

// FieldType return the type of the given field
func (c *CompositeType) FieldType(field string) string {
	f := c.Fields[field]

	if f.IsJoinedData {
		return parser.String
	}
	return f.Type
}
