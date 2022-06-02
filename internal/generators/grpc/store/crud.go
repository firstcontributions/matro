package store

const crudTmpl = `
package grpc

import (
	"context"
	"{{- .Repo -}}/internal/models/{{- .Module -}}store"
	"{{- .Repo -}}/internal/grpc/{{- plural .Module -}}/proto"

)


func (s *{{- title .Module -}}Store) Create{{- title .Name -}} (ctx context.Context, {{.Name}} *{{-  .Module -}}store. {{- title .Name}}) (* {{ .Module -}}store. {{- title .Name}}, error) {

	request := {{.Name -}}.ToProto()
	request.TimeCreated = timestamppb.Now()
	request.TimeUpdated = timestamppb.Now()

	conn, err := s.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	response, err := proto.New{{- title .Module -}}ServiceClient(conn).Create{{- title .Name -}}(ctx, request)
	if err != nil {
		return nil, err
	}
	return {{  .Module -}}store.New{{- title .Name}}().FromProto(response), nil
}

func (s *{{- title .Module -}}Store) Get{{- title .Name -}}ByID (ctx context.Context, id string) (* {{ .Module -}}store. {{- title .Name}}, error) {
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	response, err := proto.New{{- title .Module -}}ServiceClient(conn).Get{{- title .Name -}}ByID(ctx, &proto.RefByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return {{  .Module -}}store.New{{- title .Name}}().FromProto(response), nil
}
func get{{- title (plural .Name) -}}Request (
	ctx context.Context,
	ids []string,
	{{- if not (empty .SearchFields)}}
	search *string,
	{{- end}}
	{{- template "getargs" . }}
	after *string,
	before *string,
	first *int64, 
	last *int64,
) *proto.Get{{- title (plural .Name) -}}Request {
	request := &proto.Get{{- title (plural .Name) -}}Request{
		Ids: ids,
	}
	{{- if not (empty .SearchFields)}}
	if search != nil {
		request.Search = search
	}
	{{- end}}
	{{- $t := .}}
	{{- range .ReferedTypes }}
	if {{ .Name -}} != nil {
		request.{{- title .Name -}}Id = &{{- .Name -}}.Id
	}
	{{- end }}
	{{- range .Filters}}
	if {{ camel .}} != nil {
		request.{{- title . -}} = {{ camel .}}
	}
	{{- end}}
	if after != nil {
		request.After = after
	}
	if before != nil {
		request.After = before
	}
	if first != nil {
		request.First = first
	}
	if last != nil {
		request.Last = last
	}
	return request
}
func (s *{{- title .Module -}}Store) Get{{- title (plural .Name) -}} (
	ctx context.Context,
	ids []string,
	{{- if not (empty .SearchFields)}}
	search *string,
	{{- end}}
	{{- template "getargs" . }}
	after *string,
	before *string,
	first *int64, 
	last *int64,
) (
	[]*{{ .Module -}}store. {{- title .Name}}, 
	bool,
	bool,
	string,
	string,
	error,
) {
	request := get{{- title (plural .Name) -}}Request (
		ctx,
		ids,
		{{- if not (empty .SearchFields)}}
		search ,
		{{- end}}
		{{- $t := .}}
		{{- range .ReferedTypes }}
			{{ .Name -}},
		{{- end }}
		{{- range .Filters}}
			{{ camel .}},
		{{- end}}
		after,
		before,
		first, 
		last,
	)
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return nil, false, false, "", "", err
	}
	response, err := proto.New{{- title .Module -}}ServiceClient(conn).Get{{- title (plural .Name) -}}(ctx, request)
	if err != nil {
		return nil, false, false, "", "", err
	}
	
	{{plural .Name}} := []*{{  .Module -}}store.{{- title .Name -}}{}
	for _, d := range response.Data {
		{{ plural .Name}} = append({{- plural .Name}}, {{  .Module -}}store.New{{- title .Name}}().FromProto(d))
	}

	return {{ plural .Name}}, response.HasNext, response.HasPrevious, response.FirstCursor, response.LastCursor, nil
}
{{- if .Mutatable}}
func (s *{{- title .Module -}}Store) Update{{- title .Name -}} (ctx context.Context, id string, {{.Name -}}Update *{{-  .Module -}}store. {{- title .Name -}}Update) (error) {
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return err
	}
	input := {{.Name -}}Update.ToProto()
	input.Id = id

	_, err = proto.New{{- title .Module -}}ServiceClient(conn).Update{{- title .Name -}}(ctx, input)
	return err
}
{{- end}}
func (s *{{- title .Module -}}Store) Delete{{- title .Name -}}ByID (ctx context.Context, id string) (error) {
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return err
	}

	_, err = proto.New{{- title .Module -}}ServiceClient(conn).Delete{{- title .Name -}}(ctx, &proto.RefByIDRequest{Id: id})
	return err
}

{{- define "getargs"}}
{{- $t := .}}
{{- range .ReferedTypes }}
	{{.Name -}} *{{.Module.Store -}}.{{- title .Name}},
{{- end }}
{{- range .Filters}}
	{{ camel .}} *{{$t.FieldType .}},
{{- end}}
{{- end}}
`
