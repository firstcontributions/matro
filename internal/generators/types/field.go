package types

import (
	"fmt"
	"strings"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/gertd/go-pluralize"
)

// Field defines the field meta data by its type, is it a list,
// is it nullable etc..
type Field struct {
	Name             string
	Type             string
	GraphQLType      string
	IsList           bool
	IsNullable       bool
	IsPaginated      bool
	MaxCount         int
	IsQuery          bool
	Args             []Field
	IsPrimitive      bool
	IsJoinedData     bool
	IsMutatable      bool
	HardcodedFilters map[string]string
	NoGraphql        bool
	ViewerRefence    bool
}

// paginationArgs are the defualt pagination arguments should be
//  there with graphql relay paginated queries

func getPaginationArgs(f *Field) []Field {
	return []Field{
		{Name: "first", Type: parser.Int},
		{Name: "last", Type: parser.Int},
		{Name: "after", Type: parser.String},
		{Name: "before", Type: parser.String},
		{Name: "sort_order", Type: parser.String, GraphQLType: "sort_order"},
		{Name: "sort_by", Type: parser.String, GraphQLType: f.Type + "_sort_by"},
	}
}

// TODO(@gokultp) clean up this function, make it more readable
// NewField returns an instance of the field
func NewField(d *parser.Definition, typesMap map[string]*parser.Type, typeDef *parser.Type, name string, isViewerReferece bool) *Field {

	if typeDef.IsPrimitive() {
		if typeDef.Type == parser.List {
			return &Field{
				Name:        name,
				Type:        typeDef.Schema,
				IsPrimitive: true,
				NoGraphql:   typeDef.NoGraphql,
				IsList:      true,
			}
		}
		return &Field{
			Name:        name,
			Type:        typeDef.Type,
			IsPrimitive: true,
			NoGraphql:   typeDef.NoGraphql,
		}
	}

	f := &Field{
		Name:             name,
		IsList:           typeDef.Type == parser.List,
		IsPaginated:      typeDef.Paginated,
		IsJoinedData:     typeDef.JoinedData,
		HardcodedFilters: typeDef.HardcodedFilters,
		ViewerRefence:    isViewerReferece,
		MaxCount:         typeDef.MaxCount,
	}
	if typeDef.Schema == "" {
		f.Type = typeDef.Name
		f.NoGraphql = typeDef.NoGraphql
	} else if !IsCompositeType(typeDef.Schema) {
		f.Type = typeDef.Schema
		f.NoGraphql = typeDef.NoGraphql
	} else {
		f.Type = typesMap[typeDef.Schema].Name
		f.NoGraphql = typesMap[typeDef.Schema].NoGraphql
	}
	if _, ok := typesMap[f.Type]; ok {
		f.Args = getArgs(typesMap, typesMap[f.Type])
	}

	if f.IsPaginated {
		f.Args = append(f.Args, getPaginationArgs(f)...)
		f.IsQuery = true
	}
	return f
}

// getArgs gets argumets for query
func getArgs(typesMap map[string]*parser.Type, typeDef *parser.Type) []Field {
	args := []Field{}
	for _, a := range typeDef.Meta.Filters {

		for pName, pType := range typesMap[typeDef.Name].Properties {
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
						IsPrimitive: false,
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

// GoName return the field name to be used in go code
func (f *Field) GoInputName(allExported ...bool) string {
	if f.Name == "id" {
		return "Id"
	}
	return utils.ToTitleCase(f.Name)
}

// GoType return the gotype to be used in go code
// args[0] graphql enabled
// args[1] update type
func (f *Field) GoType(args ...bool) string {
	var t string
	if f.IsJoinedData {
		t = "string"
	} else {
		t = GetGoType(f.Type)
		if len(args) > 0 && args[0] {
			t = GetGoGraphQLType(f.Type)
		}
		if f.IsList {
			if IsCompositeType(f.Type) {
				t = "[]" + t
			} else {
				t = "[]*" + t
			}

		}
	}
	if (f.IsPrimitive || f.Type == "time") && len(args) > 1 && args[1] {
		t = "*" + t
	}

	return t
}

// GraphQLFormattedName returns the formatted graphql name for the field
// if it is queiriable it formats like field(args...):Type!
func (f *Field) GraphQLFormattedName() string {
	name := utils.ToCamelCase(f.Name)
	if !f.IsQuery || len(f.Args) == 0 {
		return name
	}
	args := []string{}
	for _, a := range f.Args {
		if _, ok := f.HardcodedFilters[a.Name]; ok {
			// no need to add hardcoded filters in graphql query args
			continue
		}
		args = append(args, fmt.Sprintf("%s: %s", utils.ToCamelCase(a.Name), GetGraphQLType(&a)))
	}
	return fmt.Sprintf("%s(%s)", name, strings.Join(args, ", "))
}

// GraphQLFortmattedType return the graphql type name
func (f *Field) GraphQLFortmattedType(args ...bool) string {
	var forceFieldsTobeOptional bool
	if len(args) > 0 && args[0] {
		forceFieldsTobeOptional = true
	}
	t := GetGraphQLType(f)
	if f.IsPaginated {
		plType := pluralize.NewClient().Plural(f.Type)
		t = fmt.Sprintf("%sConnection", utils.ToTitleCase(plType))
	}
	if f.IsList && !f.IsPaginated && IsCompositeType(f.Type) {
		t = fmt.Sprintf("[%s]", t)
	}
	if !forceFieldsTobeOptional && !f.IsNullable {
		t = t + "!"
	}
	return t
}

// TODO: @gokul clean this up
// GraphQLFortmattedType return the graphql type name
func (f *Field) GraphQLFortmattedInputType(args ...bool) string {
	var forceFieldsTobeOptional bool
	if len(args) > 0 && args[0] {
		forceFieldsTobeOptional = true
	}
	t := GetGraphQLType(f)
	if f.IsPaginated {
		plType := pluralize.NewClient().Plural(f.Type)
		t = fmt.Sprintf("%sConnection", utils.ToTitleCase(plType))
	}
	if f.IsJoinedData {
		return "String!"
	}
	if f.IsList && !f.IsPaginated && IsCompositeType(f.Type) {
		t = fmt.Sprintf("[%s]", t)
	}
	t = t + "Input"
	if !forceFieldsTobeOptional && !f.IsNullable {
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
