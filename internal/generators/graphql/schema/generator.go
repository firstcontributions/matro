package schema

import (
	"fmt"
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/graphql/types"
	"github.com/firstcontributions/matro/internal/generators/utils"
)

type Generator struct {
	*types.TypeDefs
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
