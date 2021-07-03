package parser

type Module struct {
	Name       string   `json:"name"`
	DataSource string   `json:"data_source"`
	DB         string   `json:"db"`
	Entities   []string `json:"entities"`
}
