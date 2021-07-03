package parser

import (
	"encoding/json"
	"errors"
	"strings"
)

const (
	String = "string"
	Int    = "int"
	Float  = "float"
	ID     = "id"
	Time   = "time"
	Bool   = "bool"
	Object = "object"
	List   = "list"
)

type Ops struct {
	Create bool
	Read   bool
	Update bool
	Delete bool
}
type Meta struct {
	SearchFields []string `json:"search_fields"`
	Filters      []string `json:"filters"`
	GraphqlOps   *Ops     `json:"graphql_ops"`
}

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

type Type struct {
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Paginated  bool             `json:"paginated"`
	Schema     string           `json:"schema"`
	JoinedData bool             `json:"joined_data"`
	Properties map[string]*Type `json:"properties"`
	Meta       Meta
}

func (t *Type) UnmarshalJSON(b []byte) error {
	// type cannont be an empty data
	ln := len(b)
	if ln == 0 {
		return errors.New("Type can't be empty")
	}
	// find the first non-space character
	var i int
	for i < ln && b[i] == ' ' {
		i++
	}
	if i == ln {
		// cannot be empty string
		return errors.New("Type can't be empty")
	}
	if b[i] != '{' {
		// this should be a string if not in the form of an obejct leteral
		// remove all the double quotes (data could be in the format "int")
		t.Type = strings.Replace(string(b), "\"", "", -1)
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

func (t *Type) IsPrimitive() bool {
	return t.Type != List && t.Type != Object
}
