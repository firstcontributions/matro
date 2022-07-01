package parser

import (
	"fmt"

	"github.com/gertd/go-pluralize"
)

// Module encapsulates the module meta data
type Module struct {
	Name       string           `json:"name"`
	DataSource string           `json:"data_source"`
	DB         string           `json:"db"`
	Entities   map[string]*Type `json:"entities"`
}

func (m *Module) Store() string {
	return fmt.Sprintf("%sstore", pluralize.NewClient().Plural(m.Name))
}

func (m *Module) Validate() error {
	return nil
}
