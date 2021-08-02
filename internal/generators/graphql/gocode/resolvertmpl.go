package gocode

var resolverTmpl = `
package schema

import (
	{{- $g := .}}
	{{- range .Modules }}
	"{{- $g.Repo -}}/internal/models/{{- .Name -}}store"
	{{- end}}
)

type Resolver struct {
	{{- range .Modules }} 
	{{ .Name -}}Store {{ .Name -}}store.Store
	{{- end}}
}

func (r *Resolver) Node(ctx context.Context, in NodeIDInput) (*NodeResolver, error) {
	id, err := ParseGraphqlID(in.ID)
	if err != nil {
		return nil, err
	}
	switch id.Type {
		{{- range .Types}}
	case "{{- .Name -}}":
		{{.Name}}Data, err := r.{{- .Module -}}Store.Get{{- title .Name -}}ByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		{{.Name -}}Node := New{{- title .Name}}({{.Name}}Data)
		return &NodeResolver{
			Node: {{.Name -}}Node,
		}
		{{- end}}
	}
	return nil, errors.New("invalid ID")
}
`
