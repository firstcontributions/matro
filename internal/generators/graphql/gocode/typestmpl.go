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
{{- template "nodeIDResolver" .}}
{{- end}}

{{- if .IsEdge}}
{{- template "edgeStruct" .}}
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

{{- define "nodeIDResolver" }}
func (n *{{ title .Name}}) ID(ctx context.Context) *graphql.ID {
	return NewIDMarshaller("{{.Name}}", *n.Id).
	ToGraphqlID()
}
{{- end}}

{{- define "edgeStruct" }}
type {{.ConnectionName}} struct {
	Edges []* {{- .EdgeName}}
	PageInfo *PageInfo
}

func New{{.ConnectionName}}(
	data []*{{- .Module -}}store.{{- title .Name}},
	hasNextPage *bool,
	hasPreviousPage *bool,
) *{{.ConnectionName}}{
	edges := []* {{- .EdgeName}}{}
	for _, d := range data {
		node := New {{- title .Name}}(d)
		edges = append(edge, &{{- .EdgeName}}{
			Node : node,
			Cursor: node.ID(),
		})
	}
	return &{{.ConnectionName}} {
		Edges: edges,
		PageInfo: &PageInfo{
			HasNextPage : hasNextPage,
			HasPreviousPage : hasPreviousPage,
		},
	}
}

type {{.EdgeName}} struct {
	Node *{{- title .Name}}
	Cursor string
}
{{- end}}

`
