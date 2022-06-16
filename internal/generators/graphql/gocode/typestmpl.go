package gocode

const typesTpl = `
package schema

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"{{- .Repo -}}/internal/models/{{- .Module.Name -}}store"
)


{{- template "typeDef" .}}
{{- template "constructor" .}}
{{- template "joinDataResolvers" .}}

{{- if (ne .Module.Name "queries")}}
{{- if .IsNode}}
{{- template "inputTypes" .}}
{{- else}}
{{- template "modelTypeAdatper" .}}
{{- end}}
{{- end}}

{{- if .IsNode }}
{{- template "nodeIDResolver" .}}

{{- if .IsEdge}}
{{- template "edgeStruct" .}}
{{- end}}
{{- end}}


{{- define "typeDef" }}
type {{ title .Name}} struct {
	ref *{{- .Module.Store -}}.{{- title .Name}}
	{{- range .Fields}}
	{{- template "fieldDef" .}}
	{{- end}}
}
{{- end}}

{{- define "inputTypes" }}
type Create{{- title .Name -}}Input struct {
	{{- range .Fields}}
	{{- if (not (or (isAditField .Name) .IsQuery))}}
	{{- template "inputFieldDef" .}}
	{{- end}}
	{{- end}}

	{{- range .ReferedTypes}}
	{{title .Name}}ID graphql.ID
	{{- end}}
}

func (n *Create{{- title .Name -}}Input) ToModel() (*{{- .Module.Store -}}.{{- title .Name}}, error){
	if n == nil {
		return nil, nil
	}
	{{- range .ReferedTypes}}
	{{- if (not .IsViewerType)}}
	{{ .Name}}ID, err := ParseGraphqlID(n.{{- title .Name}}ID)
	if err != nil {
		return nil, err
	}
	{{- end}}
	{{- end}}
	
	return &{{- .Module.Store -}}.{{- title .Name}} {
		{{- range .Fields }}
		{{- if (not (or (isAditField .Name) .IsQuery))}}
			{{- template "inputTypeField" .}}
		{{- end}}
		{{- end}}

		{{- range .ReferedTypes}}
		{{- if (not .IsViewerType)}}
		{{title .Name}}ID : {{- .Name}}ID.ID,
		{{- end}}
		{{- end}}
	}, nil
}

{{- if (and .IsNode .GraphqlOps.Update)}}
type Update{{- title .Name -}}Input struct {
	ID graphql.ID
	{{- range .Fields}}
	{{- if (and .IsMutatable (not (or (isAditField .Name) .IsQuery .NoGraphql)))}}
	{{- if  (not (and .IsJoinedData  .IsList))}}
	{{.GoName}} {{.GoType true true}}
	{{- end}}
	{{- end}}
	{{- end}}
}

func (n *Update{{- title .Name -}}Input) ToModel() *{{- .Module.Store -}}.{{- title .Name -}}Update {
	if n == nil {
		return nil
	}
	 return &{{- .Module.Store -}}.{{- title .Name -}}Update {
		{{- range .Fields }}
		{{- if (and .IsMutatable (not (or (isAditField .Name) .IsQuery)))}}
			{{- template "inputTypeField" .}}
		{{- end}}
		{{- end}}
	}
}
{{- end }}
{{- end}}

{{- define "inputTypeField" }}
{{- if  (not (or (and .IsJoinedData  .IsList) .NoGraphql .ViewerRefence))}}
	{{- if (and (not .IsJoinedData) (isCompositeType .Type))}}
		{{.GoInputName}} :n.{{- .GoName }}.ToModel(),
	{{- else}}
		{{- if (eq .Type "int")}}
			{{.GoInputName}} : int64(n.{{- .GoName }}), 
		{{- else }}
			{{- if (eq .Type "time")}}
				{{.GoInputName}} : n.{{- .GoName }}.Time, 
			{{- else}}
				{{.GoInputName}} : n.{{- .GoInputName }}, 
			{{- end}}
		{{- end}}
	{{- end}}
{{- end}}
{{- end}}

{{- define "fieldDef" }}
	{{- if  (not (or (and .IsJoinedData  .IsList) .NoGraphql))}}
	{{.GoName}} {{.GoType true}}
	{{- end}}
{{- end}}
{{- define "inputFieldDef" }}
	{{- if  (not (or (and .IsJoinedData  .IsList) .NoGraphql .ViewerRefence))}}
	{{.GoInputName}} {{.GoType true}}
	{{- end}}
{{- end}}

{{- define "constructor" }}
{{- if .AllReferedFields}}
func New {{- title .Name}} () *{{- title .Name}} {
{{- else}}
func New {{- title .Name}} (m *{{.Module.Name -}}store.{{-  title .Name}}) *{{- title .Name}} {
	if m == nil {
		return nil
	}
{{- end}}
	return &{{- title .Name}} {
		ref : m,
		{{- range .Fields}}
		{{- if  (not (or (and .IsJoinedData  .IsList) .NoGraphql))}}
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

	data, err := storemanager.FromContext(ctx).{{- title (plural $returntype.Module.Name) }}Store.Get{{- title $returntype.Name -}}ByID(ctx, n.{{- .GoName}})
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
	data []*{{- .Module.Name -}}store.{{- title .Name}},
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
			{{- if (eq .Module.DB "")}}
			Cursor: d.Cursor,
			{{- else}}
			Cursor: cursor.NewCursor(d.Id, d.TimeCreated).String(),
			{{- end}}
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


{{- define "modelTypeAdatper" }}
func (n *{{title .Name}}) ToModel() *{{- .Module.Store -}}.{{- title .Name}}{
	if n == nil {
		return nil
	}
	 return &{{- .Module.Store -}}.{{- title .Name}} {
		{{- range .Fields }}
		{{- if  (not (or (and .IsJoinedData  .IsList) .NoGraphql))}}
		{{- if (and (not .IsJoinedData) (isCompositeType .Type))}}
		{{.GoInputName}} :n.{{- .GoName }}.ToModel(),
		{{- else}}
		{{- if (eq .Type "int")}}
		{{.GoInputName}} : int64(n.{{- .GoName }}), 
		{{- else }}
		{{- if (eq .Type "time")}}
		{{.GoInputName}} : n.{{- .GoName }}.Time, 
		{{- else}}
		{{.GoInputName}} : n.{{- .GoInputName }}, 
		{{- end}}
		{{- end}}
		{{- end}}
		{{- end}}

		{{- end}}
	}
}
{{- end}}

`
