package types

import (
	"fmt"
	"strings"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

// CompositeType defines a non trivial type with a set
// of informations like what all fields, opeations supported, etc.
type CompositeType struct {
	Name          string
	Fields        map[string]*Field
	ReferedFields []string
	IsNode        bool
	IsEdge        bool
	Filters       []string
	GraphqlOps    *parser.Ops
	SearchFields  []string
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
	}
	if f.IsPaginated {
		f.Args = append(f.Args, paginationArgs...)
		f.IsQuery = true
	}
	return f
}

// GoName return the field name to be used in go code
func (f *Field) GoName() string {
	if f.Name == "id" {
		return "Id"
	}
	if f.IsJoinedData {
		return f.Name
	}
	return utils.ToTitleCase(f.Name)
}

// GoType return the gotype to be used in go code
func (f *Field) GoType(graphqlEnabled ...bool) string {
	if f.IsJoinedData {
		return "*string"
	}
	t := getGoType(f.Type)
	if len(graphqlEnabled) > 0 && graphqlEnabled[0] {
		t = getGoGraphQLType(f.Type)
	}
	t = "*" + t
	if f.IsList {
		t = "[]" + t
	}
	return t
}

// GraphQLFormattedName returns the formatted graphql name for the field
// if it is queiriable it formats like field(args...):Type!
func (f *Field) GraphQLFormattedName() string {
	if !f.IsQuery {
		return f.Name
	}
	args := []string{}
	for _, a := range f.Args {
		args = append(args, fmt.Sprintf("%s: %s", a.Name, getGraphQLType(a.Type)))
	}
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
}

// GraphQLFortmattedType return the graphql type name
func (f *Field) GraphQLFortmattedType() string {
	t := getGraphQLType(f.Type)
	if f.IsPaginated {
		t = fmt.Sprintf("%ssConnection", utils.ToTitleCase(f.Type))
	}
	if f.IsList {
		t = fmt.Sprintf("[%s]", t)
	}
	if !f.IsNullable {
		t = t + "!"
	}
	return t
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
	return &CompositeType{
		Name:         typeDef.Name,
		Fields:       fields,
		IsNode:       isNode,
		GraphqlOps:   typeDef.Meta.GraphqlOps,
		SearchFields: typeDef.Meta.SearchFields,
		Filters:      typeDef.Meta.Filters,
	}
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
	return fmt.Sprintf("%ssConnection", utils.ToTitleCase(c.Name))
}

// FieldType return the type of the given field
func (c *CompositeType) FieldType(field string) string {
	f := c.Fields[field]
	if f.IsJoinedData {
		return parser.String
	}
	return f.Type
}
