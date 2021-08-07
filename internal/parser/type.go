package parser

import (
	"encoding/json"
	"errors"
	"strings"
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

var (
	errEmptyType = errors.New("Type can't be empty")
)

// Meta defines the data schema meta data properties
type Meta struct {
	SearchFields    []string `json:"search_fields"`
	Filters         []string `json:"filters"`
	GraphqlOps      *Ops     `json:"graphql_ops"`
	MutatableFields []string `json:"mutatable_fields"`
}

// Ops defines the supported graphql CRUD operations
type Ops struct {
	Create bool
	Read   bool
	Update bool
	Delete bool
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
			m.Create = true
		case 'R':
			m.Read = true
		case 'U':
			m.Update = true
		case 'D':
			m.Delete = true
		}
	}
	return nil
}

// Type encapsulates the type definition meta data
type Type struct {
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Paginated  bool             `json:"paginated"`
	Schema     string           `json:"schema"`
	JoinedData bool             `json:"joined_data"`
	Properties map[string]*Type `json:"properties"`
	Meta       Meta
}

// validateTypeAndGetFirstNonEmptyIdx validates the given type string binary
func validateTypeAndGetFirstNonEmptyIdx(b []byte) (int, error) {
	// type cannont be an empty data
	ln := len(b)
	if ln == 0 {
		return 0, errEmptyType
	}
	// find the first non-space character
	var i int
	for i < ln && b[i] == ' ' {
		i++
	}
	if i == ln {
		// cannot be empty string
		return 0, errEmptyType
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
		// this should be a string if not in the form of an obejct leteral
		// remove all the double quotes (data could be in the format "int")
		t.Type = strings.ReplaceAll(string(b), "\"", "")
		return nil
	}
	// if object literal, use normal json unmarshal func for the struct
	// we cannot call json.Unmarshal on this struct as it will be an recursive call
	// to the same func.

	// declaring var data with inline struct type
	data := struct {
		Name       string           `json:"name"`
		Type       string           `json:"type"`
		Paginated  bool             `json:"paginated"`
		Schema     string           `json:"schema"`
		JoinedData bool             `json:"joined_data"`
		Properties map[string]*Type `json:"properties"`
		Meta       Meta             `json:"meta"`
	}{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	t.Type = data.Type
	t.Paginated = data.Paginated
	t.JoinedData = data.JoinedData
	t.Schema = data.Schema
	t.Properties = data.Properties
	t.Name = data.Name
	t.Meta = data.Meta
	return nil
}

// IsPrimitive says the type is primitive or not
func (t *Type) IsPrimitive() bool {
	return t.Type != List && t.Type != Object
}
