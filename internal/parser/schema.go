package parser

import (
	"encoding/json"
	"io/ioutil"
)

// Definition encapsulates the input schema definiton
type Definition struct {
	DataSchema map[string]*Type `json:"data_schema"`
	Modules    []Module         `json:"modules"`
	Repo       string           `json:"repo"`
	Queries    []*Type          `json:"high_level_queries"`
}

// NewDefinition return an instance of Definition
func NewDefinition() *Definition {
	return &Definition{
		DataSchema: map[string]*Type{},
	}
}

// ParseFromFile parses the definitions from the given json file
// returns an instance of the same object and error if any.
//
// the duplicate return of object helps in chaining the functions like
// NewDefinition().ParseFromFile
func (d *Definition) ParseFromFile(path string) (*Definition, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return d, err
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return d, err
	}
	return d, nil
}
