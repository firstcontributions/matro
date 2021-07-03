package schema

import (
	"fmt"
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

type Generator struct {
	*types.TypeDefs
	Path string
}

func NewGenerator(path string, d *parser.Definition) *Generator {
	td := types.NewTypeDefs(path, d)
	return &Generator{
		Path:     path,
		TypeDefs: td,
	}
}
func (g *Generator) Generate() error {
	path := fmt.Sprintf("%s/assets", g.Path)
	t, err := template.New("graphql").
		Funcs(g.FuncMap()).
		Parse(schemaTmpl)
	if err != nil {
		return err
	}
	fw, err := utils.GetFileWriter(path, "schema.graphql")
	if err != nil {
		return err
	}
	defer fw.Close()
	return t.Execute(fw, g)
}
