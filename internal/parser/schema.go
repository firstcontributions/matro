package parser

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Defaults struct {
	ViewerType string `json:"viewer_type"`
}

// Definition encapsulates the input schema definiton
type Definition struct {
	Modules  []Module  `json:"modules"`
	Repo     string    `json:"repo"`
	Queries  []*Type   `json:"high_level_queries"`
	Defaults *Defaults `json:"defaults"`
}

// NewDefinition return an instance of Definition
func NewDefinition() *Definition {
	return &Definition{}
}

// ParseFrom parses the definitions from the given reader
// returns an instance of the same object and error if any.
// the duplicate return of object helps in chaining the functions like
// NewDefinition().ParseFrom()
func (d *Definition) ParseFrom(reader io.Reader) (*Definition, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return d, err
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return d, err
	}
	return d, nil
}
