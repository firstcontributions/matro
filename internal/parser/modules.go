package parser

// Module encapsulates the module meta data
type Module struct {
	Name       string           `json:"name"`
	DataSource string           `json:"data_source"`
	DB         string           `json:"db"`
	Entities   map[string]*Type `json:"entities"`
}
