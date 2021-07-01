package gocode

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"

	"github.com/firstcontributions/matro/internal/generators/graphql/types"
	"github.com/firstcontributions/matro/internal/generators/utils"
)

type Generator struct {
	*types.TypeDefs
}

func (g *Generator) Generate() error {
	path := fmt.Sprintf("%s/intern/graphql/schema", g.Path)
	for _, t := range g.Types {
		if err := g.generateCodeFromTemplate(typesTpl, t, path, t.Name+".go"); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) generateCodeFromTemplate(tmpl string, data interface{}, path, filename string) error {
	fw, err := utils.GetFileWriter(path, filename)
	defer fw.Close()
	if err != nil {
		return err
	}
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
	code := b.Bytes()
	code, err = format.Source(code)
	if err != nil {
		return err
	}
	_, err = fw.Write(code)
	return err
}
