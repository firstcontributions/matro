package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/firstcontributions/matro/internal/merrors"
)

const (
	// String defines the string type
	String = "string"
	// Int defines the int type supported, here the system uses int32 internally
	// as google v8 engine and graphql supports only int32
	Int = "int"
	// Float defines the float type
	Float = "float"
	// ID defines the ID type, the sustem uses string version of uuid
	ID = "id"
	// Time defines timestamp type, internally uses time.Time
	Time = "time"
	// Bool defines boolean type
	Bool = "bool"
	// Object defines object type schema for composite types
	Object = "object"
	// List defines list type schema for composite types
	List = "list"
)

// Meta defines the data schema meta data properties
type Meta struct {
	SearchFields       []string `json:"search_fields"`
	Filters            []string `json:"filters"`
	GraphqlOps         Ops      `json:"graphql_ops"`
	MutatableFields    []string `json:"mutatable_fields"`
	SortBy             []string `json:"sort_by"`
	ViewerRefenceField string   `json:"viewer_reference_field"`
}

// Ops defines the supported graphql CRUD operations
type Ops struct {
	create bool
	read   bool
	update bool
	delete bool
}

func (o *Ops) Create() bool {
	return o != nil && o.create
}
func (o *Ops) Update() bool {
	return o != nil && o.update
}
func (o *Ops) Read() bool {
	return o != nil && o.read
}
func (o *Ops) Delete() bool {
	return o != nil && o.delete
}
func (o *Ops) Union(op2 Ops) {
	o.create = o.create || op2.create
	o.read = o.read || op2.read
	o.update = o.update || op2.update
	o.delete = o.delete || op2.delete
}

// UnmarshalJSON is a custom json unmarshaller
// the operations value will be given as a string `CRUD`
// for eg: to support only Create and Read ops, it will be `CR`
// this custom unmarshaller will parse the string and construct
// the Ops object
func (m *Ops) UnmarshalJSON(data []byte) error {
	for _, b := range data {
		switch b {
		case 'C':
			m.create = true
		case 'R':
			m.read = true
		case 'U':
			m.update = true
		case 'D':
			m.delete = true
		}
	}
	return nil
}

// Type encapsulates the type definition meta data
type Type struct {
	_Type
}

type _Type struct {
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	Paginated        bool              `json:"paginated"`
	Schema           string            `json:"schema"`
	JoinedData       bool              `json:"joined_data"`
	Properties       map[string]*Type  `json:"properties"`
	Meta             Meta              `json:"meta"`
	HardcodedFilters map[string]string `json:"hardcoded_filters"`
	NoGraphql        bool              `json:"no_graphql"`
}

// validateTypeAndGetFirstNonEmptyIdx validates the given type string binary
func validateTypeAndGetFirstNonEmptyIdx(b []byte) (int, error) {
	// type cannont be an empty data
	ln := len(b)
	if ln == 0 {
		return 0, merrors.ErrEmptyType
	}
	// find the first non-space character
	var i int
	for i < ln && b[i] == ' ' {
		i++
	}
	if i == ln {
		// cannot be empty string
		return 0, merrors.ErrEmptyType
	}
	return i, nil
}

// UnmarshalJSON is a custom json unmarshaller
// type can be either a sting or an object in a given format
// eg: "int" or {type: "int", meta: {}}
func (t *Type) UnmarshalJSON(b []byte) error {
	i, err := validateTypeAndGetFirstNonEmptyIdx(b)
	if err != nil {
		return err
	}
	if b[i] != '{' {
		// this should be a string if not in the form of an obejct literal
		// remove all the double quotes (data could be in the format "string")
		t.Type = strings.ReplaceAll(string(b), "\"", "")
		return nil
	}
	// if object literal, use normal json unmarshal func for the struct
	// we cannot call json.Unmarshal on this struct as it will be an recursive call
	// to the same func.

	if err := json.Unmarshal(b, &t._Type); err != nil {
		return fmt.Errorf("error on unmarshalling, %w [%v]", merrors.ErrInvalidTypeDefinition, err)
	}
	return nil
}

// IsPrimitive says the type is primitive or not
func (t *Type) IsPrimitive() bool {
	switch t.Type {
	case Int, Bool, Float, ID, String, Time:
		return true
	}
	return false
}

// Validate will validate parsed type
func (t *Type) Validate() error {
	if t.Type == "" {
		return merrors.ErrEmptyType
	}
	if (t.Type == List || t.Type == Object) && t.Schema == "" && t.Properties == nil {
		return merrors.ErrNoSchemaDefinedForCustomType
	}
	if (t.Type == List || t.Type == Object) && t.Properties != nil && t.Name == "" {
		return merrors.ErrNoNameForInlineCustomTypes
	}
	for _, f := range t.Meta.SearchFields {
		if _, found := t.Properties[f]; !found {
			return fmt.Errorf("%w field: %v", merrors.ErrUndefinedSearchField, f)
		}
	}
	if err := t.raiseErrorIfFieldsNotDefined(t.Meta.SearchFields, merrors.ErrUndefinedSearchField); err != nil {
		return err
	}
	if err := t.raiseErrorIfFieldsNotDefined(t.Meta.MutatableFields, merrors.ErrUndefinedMutableField); err != nil {
		return err
	}
	if err := t.raiseErrorIfFieldsNotDefined(t.Meta.Filters, merrors.ErrUndefinedFilter); err != nil {
		return err
	}
	if err := t.raiseErrorIfFieldsNotDefined(t.Meta.SortBy, merrors.ErrUndefinedSortField); err != nil {
		return err
	}
	return nil
}

func (t *Type) raiseErrorIfFieldsNotDefined(fields []string, err error) error {
	for _, f := range fields {
		if _, found := t.Properties[f]; !found {
			return fmt.Errorf("field: %v [%w]", f, err)
		}
	}
	return nil
}
