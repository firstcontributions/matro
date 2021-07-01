package parser

import (
	"encoding/json"
	"io/ioutil"
)

type Definition struct {
	DataSchema map[string]*Type `json:"data_schema"`
}

func NewDefinition() *Definition {
	return &Definition{
		DataSchema: map[string]*Type{},
	}
}

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
