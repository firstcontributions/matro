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
	{{- if (eq .Module.DataSource "external_apis")}}
	return NewIDMarshaller(NodeType{{- title .Name}}, n.Id, false).
	{{- else}}
	return NewIDMarshaller(NodeType{{- title .Name}}, n.Id, true).
	{{- end}}
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
	filters *{{- .Module.Name -}}store.{{- title .Name -}}Filters
}

func New{{.ConnectionName}}(
	filters *{{- .Module.Store -}}.{{- title .Name -}}Filters,
	data []*{{- .Module.Store -}}.{{- title .Name}},
	hasNextPage bool,
	hasPreviousPage bool,
	cursors []string, 
) *{{.ConnectionName}}{
	edges := []* {{- .EdgeName}}{}
	for i, d := range data {
		node := New {{- title .Name}}(d)

		edges = append(edges, &{{- .EdgeName}}{
			Node : node,
			Cursor: cursors[i],
		})
	}
	var startCursor, endCursor *string
	if len(cursors) > 0 {
		startCursor = &cursors[0]
		endCursor = &cursors[len(cursors)-1]
	}
	return &{{.ConnectionName}} {
		filters: filters,
		Edges: edges,
		PageInfo: &PageInfo{
			HasNextPage : hasNextPage,
			HasPreviousPage : hasPreviousPage,
			StartCursor :startCursor,
			EndCursor :endCursor,
		},
	}
}

func (c {{.ConnectionName}}) TotalCount (ctx context.Context) (int32, error) {
	count, err := storemanager.FromContext(ctx).{{- title (plural .Module.Name) }}Store.Count{{- plural (title .Name) -}}(ctx, c.filters)
	return int32(count), err
}

{{- if (ne .ViewerRefenceField "") }}
func (c {{.ConnectionName}}) HasViewerAssociation (ctx context.Context) (bool, error) {
	session := session.FromContext(ctx)
	if session == nil {
		return false, errors.New("Unauthorized")
	}
	userID := session.UserID()

	newFilter := *c.filters
	newFilter.{{- title .ViewerRefenceField}} = &userID

	data, err := storemanager.FromContext(ctx).{{- title (plural .Module.Name) }}Store.GetOne{{- title .Name -}}(ctx, c.filters)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}
{{- end}}

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
