package gocode

const typesTpl = `
package schema

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"{{- .Repo -}}/internal/models/{{- .Module -}}store"
)


{{- template "typeDef" .}}
{{- template "constructor" .}}

{{- if .IsNode }}
func (n *{{ title .Name}}) ID(ctx context.Context) *graphql.ID {
	return NewIDMarshaller("{{.Name}}", *n.Id).
	ToGraphqlID()
}
{{- end}}

{{- define "typeDef" }}
type {{ title .Name}} struct {
	{{- range .Fields}}
	{{- template "fieldDef" .}}
	{{- end}}
}
{{- end}}


{{- define "fieldDef" }}
	{{- if  (not (and .IsJoinedData  .IsList))}}
	{{.GoName}} {{.GoType true}}
	{{- end}}
{{- end}}

{{- define "constructor" }}
func New {{- title .Name}} (m *{{.Module -}}store.{{-  title .Name}}) *{{- title .Name}} {
	return &{{- title .Name}} {
		{{- range .Fields}}
		{{- if  (not (and .IsJoinedData  .IsList))}}
		{{.GoName}} : m.{{- .GoName true}}, 
		{{- end}}
		{{- end}}
	}
}
{{- end}}

`
