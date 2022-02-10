package gocode

var queryResolverTmpl = `
package schema

type {{title .Query.Name -}}Input struct {
	{{- range .Query.Args}}
	{{.GoName true}} *{{.GoType true}}
	{{- end}}
}
{{- if .Query.Parent }}
func (n *{{- title .Query.Parent.Name -}}) {{title .Query.Name}}(ctx context.Context, in *{{title .Query.Name -}}Input) (*{{.ReturnType.ConnectionName}}, error) {
	store := storemanager.FromContext(ctx)
{{- else }}
func (r *Resolver) {{title .Query.Name}}(ctx context.Context, in *{{title .Query.Name -}}Input) (*{{.ReturnType.ConnectionName}}, error) {
{{- end}}
	var first, last *int64
	if in.First != nil {
		tmp := int64(*in.First)
		first = &tmp
	}
	if in.Last != nil {
		tmp := int64(*in.Last)
		last = &tmp
	}
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err :=  store.{{- title .ReturnType.Module -}}Store.Get{{- plural (title .ReturnType.Name)}} (
		ctx,
		{{- template "getargs" .}}
	)
	if err != nil {
		return nil, err
	}
	return New{{- .ReturnType.ConnectionName}}(data, hasNextPage, hasPreviousPage, &firstCursor, &lastCursor), nil
}


{{- define "getargs"}}
{{- $q := .Query}}
{{- if (isElemOfStrArray ($q.ArgNames) "ids")}}
	in.Ids,
{{- else}}
	nil,
{{- end}}
{{- if not (empty .ReturnType.SearchFields)}}
	nil,
{{- end}}
{{- range .ReturnType.ReferedFields }}
		&n.Id,
{{- end }}
{{- range .ReturnType.Filters}}
	{{- if (isElemOfStrArray ($q.ArgNames) .)}}
		in.{{- title .}},
	{{- else}}
		nil,
	{{- end}}
{{- end}}
	in.After,
	in.Before,
	first,
	last, 
{{- end}}

`
