package types

import (
	"fmt"
	"strings"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

type CompositeType struct {
	Name   string
	Fields []*Field
	IsNode bool
	IsEdge bool
}
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

var paginationArgs = []Field{
	{Name: "first", Type: parser.Int},
	{Name: "last", Type: parser.Int},
	{Name: "after", Type: parser.String},
	{Name: "before", Type: parser.String},
}

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

func (f *Field) GoName() string {
	if f.Name == "id" {
		return f.Name
	}
	return utils.ToTitleCase(f.Name)
}

func (f *Field) GoType() string {
	t := getGoType(f.Type)
	t = "*" + t
	if f.IsList {
		t = "[]" + t
	}
	return t
}
func (f *Field) FormattedName() string {
	if !f.IsQuery {
		return f.Name
	}
	args := []string{}
	for _, a := range f.Args {
		args = append(args, fmt.Sprintf("%s: %s", a.Name, getGraphQLType(a.Type)))
	}
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
}

func (f *Field) FortmattedType() string {
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
func NewCompositeType(d *parser.Definition, typeDef *parser.Type) *CompositeType {
	fields := []*Field{}
	isNode := false
	for field, def := range typeDef.Properties {
		if field == "id" {
			isNode = true
		}
		fields = append(fields, NewField(d, def, field))
	}
	return &CompositeType{
		Name:   typeDef.Name,
		Fields: fields,
		IsNode: isNode,
	}
}

func (c *CompositeType) EdgeTypes() *utils.Set {
	s := utils.NewSet()
	for _, f := range c.Fields {
		if f.IsPaginated && f.IsList {
			s.Add(f.Type)
		}
	}
	return s
}
func (c *CompositeType) EdgeType() string {
	return fmt.Sprintf("%sEdge", utils.ToTitleCase(c.Name))
}

func (c *CompositeType) ConnType() string {
	return fmt.Sprintf("%ssConnection", utils.ToTitleCase(c.Name))
}
