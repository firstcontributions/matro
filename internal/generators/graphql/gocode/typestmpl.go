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
{{- template "joinDataResolvers" .}}

{{- if .IsNode }}
{{- template "nodeIDResolver" .}}

{{- if .IsEdge}}
{{- template "edgeStruct" .}}
{{- end}}
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
		{{- if (and (not .IsJoinedData) (isCompositeType .Type))}}
		{{.GoName}} : New{{- title .Type -}}(m.{{- .GoName true}}),
		{{- else}}
		{{- if (eq .Type "int")}}
		{{.GoName}} : int32(m.{{- .GoName true}}), 
		{{- else }}
		{{- if (eq .Type "time")}}
		{{.GoName}} : graphql.Time{Time: m.{{- .GoName true}}}, 
		{{- else}}
		{{.GoName}} : m.{{- .GoName true}}, 
		{{- end}}
		{{- end}}
		{{- end}}
		{{- end}}
		{{- end}}
	}
}
{{- end}}

{{- define "nodeIDResolver" }}
func (n *{{ title .Name}}) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller("{{.Name}}", n.Id).
	ToGraphqlID()
}
{{- end}}

{{- define "joinDataResolvers" }}
{{- $t := .}}
{{- range .Fields}}
{{- if  (and .IsJoinedData  (not .IsList))}}
{{- $returntype := (getTypeFromMap $t.Types .Type )}}
func (n *{{ title $t.Name}}) {{title .GoName}} (ctx context.Context) (*{{- title $returntype.Name}}, error) {
	store := ctx.Value("store").(*Store)

	data, err := store.{{- plural $returntype.Module }}Store.Get{{- title $returntype.Name -}}ByID(ctx, n.{{- .GoName}})
	if err != nil {
		return nil, err
	}
	return New{{- title $returntype.Name}}(data), nil
}
{{- end}}
{{- end}}
{{- end}}


{{- define "edgeStruct" }}
type {{.ConnectionName}} struct {
	Edges []* {{- .EdgeName}}
	PageInfo *PageInfo
}


func New{{.ConnectionName}}(
	data []*{{- .Module -}}store.{{- title .Name}},
	hasNextPage bool,
	hasPreviousPage bool,
	firstCursor *string, 
	lastCursor *string,
) *{{.ConnectionName}}{
	edges := []* {{- .EdgeName}}{}
	for _, d := range data {
		node := New {{- title .Name}}(d)

		edges = append(edges, &{{- .EdgeName}}{
			Node : node,
			Cursor: cursor.NewCursor(d.Id, d.TimeCreated).String(),
		})
	}
	return &{{.ConnectionName}} {
		Edges: edges,
		PageInfo: &PageInfo{
			HasNextPage : hasNextPage,
			HasPreviousPage : hasPreviousPage,
			StartCursor :firstCursor,
			EndCursor :lastCursor,
		},
	}
}

type {{.EdgeName}} struct {
	Node *{{- title .Name}}
	Cursor string
}
{{- end}}

`
