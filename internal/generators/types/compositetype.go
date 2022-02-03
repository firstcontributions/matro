package types

import (
	"fmt"
	"strings"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/gertd/go-pluralize"
)

// CompositeType defines a non trivial type with a set
// of informations like what all fields, opeations supported, etc.
type CompositeType struct {
	Name            string
	Fields          map[string]*Field
	ReferedFields   []string
	IsNode          bool
	IsEdge          bool
	Filters         []string
	GraphqlOps      *parser.Ops
	SearchFields    []string
	MutatableFields []string
	Module          string
}

// Field defines the field meta data by its type, is it a list,
// is it nullable etc..
type Field struct {
	Name         string
	Type         string
	IsList       bool
	IsNullable   bool
	IsPaginated  bool
	IsQuery      bool
	Args         []Field
	IsPrimitive  bool
	IsJoinedData bool
	IsMutatable  bool
}

// paginationArgs are the defualt pagination arguments should be
//  there with graphql relay paginated queries
var paginationArgs = []Field{
	{Name: "first", Type: parser.Int},
	{Name: "last", Type: parser.Int},
	{Name: "after", Type: parser.String},
	{Name: "before", Type: parser.String},
}

// NewField returns an instance of the field
func NewField(d *parser.Definition, typeDef *parser.Type, name string) *Field {
	if typeDef.IsPrimitive() {
		return &Field{
			Name:        name,
			Type:        typeDef.Type,
			IsPrimitive: true,
		}
	}
	f := &Field{
		Name:         name,
		Type:         d.DataSchema[typeDef.Schema].Name,
		IsList:       typeDef.Type == parser.List,
		IsPaginated:  typeDef.Paginated,
		IsJoinedData: typeDef.JoinedData,
		Args:         getArgs(d, typeDef),
	}
	if f.IsPaginated {
		f.Args = append(f.Args, paginationArgs...)
		f.IsQuery = true
	}
	return f
}

// getArgs gets argumets for query
func getArgs(d *parser.Definition, typeDef *parser.Type) []Field {
	args := []Field{}
	for _, a := range typeDef.Meta.Filters {

		for pName, pType := range d.DataSchema[typeDef.Schema].Properties {
			if pName == a {
				if pType.IsPrimitive() {
					args = append(args, Field{
						Name:        a,
						Type:        pType.Type,
						IsPrimitive: true,
					})
				} else {
					args = append(args, Field{
						Name:        a,
						Type:        parser.String,
						IsPrimitive: true,
					})
				}
				break
			}
		}
	}
	return args
}

// GoName return the field name to be used in go code
func (f *Field) GoName(allExported ...bool) string {
	exported := len(allExported) > 0 && allExported[0]
	if f.Name == "id" {
		return "Id"
	}
	if !exported && f.IsJoinedData {
		return utils.ToCamelCase(f.Name)
	}
	return utils.ToTitleCase(f.Name)
}

// GoType return the gotype to be used in go code
func (f *Field) GoType(graphqlEnabled ...bool) string {
	if f.IsJoinedData {
		return "string"
	}
	t := GetGoType(f.Type)
	if len(graphqlEnabled) > 0 && graphqlEnabled[0] {
		t = GetGoGraphQLType(f.Type)
	}
	if f.IsList {
		t = "[]" + t
	}
	return t
}

// GraphQLFormattedName returns the formatted graphql name for the field
// if it is queiriable it formats like field(args...):Type!
func (f *Field) GraphQLFormattedName() string {
	name := utils.ToCamelCase(f.Name)
	if !f.IsQuery {
		return name
	}
	args := []string{}
	for _, a := range f.Args {
		args = append(args, fmt.Sprintf("%s: %s", utils.ToCamelCase(a.Name), GetGraphQLType(&a)))
	}
	return fmt.Sprintf("%s(%s)", name, strings.Join(args, ", "))
}

// GraphQLFortmattedType return the graphql type name
func (f *Field) GraphQLFortmattedType() string {
	t := GetGraphQLType(f)
	if f.IsPaginated {
		plType := pluralize.NewClient().Plural(f.Type)
		t = fmt.Sprintf("%sConnection", utils.ToTitleCase(plType))
	}
	if f.IsList && !f.IsPaginated {
		t = fmt.Sprintf("[%s]", t)
	}
	if !f.IsNullable {
		t = t + "!"
	}
	return t
}

func (f *Field) ArgNames() []string {
	args := []string{}
	for _, a := range f.Args {
		args = append(args, a.Name)
	}
	return args
}

var auditFields = []*Field{
	{
		Name: "time_created",
		Type: "time",
	},
	{
		Name:        "time_updated",
		Type:        "time",
		IsMutatable: true,
	},
}

// NewCompositeType return an instance of the CompositeType
func NewCompositeType(d *parser.Definition, typeDef *parser.Type) *CompositeType {
	fields := map[string]*Field{}
	isNode := false
	for field, def := range typeDef.Properties {
		if field == "id" {
			isNode = true
		}
		fields[field] = NewField(d, def, field)
	}
	for _, f := range auditFields {
		fields[f.Name] = f
	}
	for _, mf := range typeDef.Meta.MutatableFields {
		fields[mf].IsMutatable = true
	}
	return &CompositeType{
		Name:            typeDef.Name,
		Fields:          fields,
		IsNode:          isNode,
		GraphqlOps:      typeDef.Meta.GraphqlOps,
		SearchFields:    typeDef.Meta.SearchFields,
		Filters:         typeDef.Meta.Filters,
		MutatableFields: typeDef.Meta.MutatableFields,
	}
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
func (c *CompositeType) EdgeFields() *utils.Set {
	s := utils.NewSet()
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
