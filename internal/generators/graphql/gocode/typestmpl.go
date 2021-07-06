package gocode

const typesTpl = `
package schema

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
)


{{- template "typeDef" .}}

{{- if .IsNode }}
func (n *{{ title .Name}}) ID(ctx context.Context) *graphql.ID {
	return n.Id
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

`
