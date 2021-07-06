package gocode

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/firstcontributions/matro/internal/parser"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/generators/utils"
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
	path := fmt.Sprintf("%s/internal/graphql/schema", g.Path)
	for _, t := range g.Types {
		if err := g.generateCodeFromTemplate(typesTpl, t, path, t.Name+".go"); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) generateCodeFromTemplate(tmpl string, data interface{}, path, filename string) error {

	var b bytes.Buffer

	t, err := template.New("go").
		Funcs(g.FuncMap()).
		Parse(tmpl)
	if err != nil {
		return err
	}
	if err := t.Execute(&b, data); err != nil {
		return err
	}
	return utils.WriteCodeToGoFile(path, filename, b.Bytes())
}
