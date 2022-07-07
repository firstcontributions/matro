package gocode

var queryResolverTmpl = `
package schema

{{- if (ne (len .Query.Args) 0)}}
type {{ .Query.InputName }} struct {
	{{- $q := .Query}}
	{{- range .Query.Args}}
	{{- if (not (isHardCodedFilter $q.HardcodedFilters .Name))}}
	{{.GoName true}} *{{.GoType true}}
	{{- end}}
	{{- end}}
}
{{- end}}
{{- if .Query.IsPaginated }}
{{- if .Query.Parent }}
func (n *{{- title .Query.Parent.Name -}}) {{title .Query.Name}}(ctx context.Context, in *{{ .Query.InputName -}}) (*{{.ReturnType.ConnectionName}}, error) {
{{- else }}
func (r *Resolver) {{title .Query.Name}}(ctx context.Context, in *{{title .Query.InputName -}}) (*{{.ReturnType.ConnectionName}}, error) {
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
	store := storemanager.FromContext(ctx)
	{{- $q := .Query}}
	{{- range .Query.Args}}
		{{- if (isHardCodedFilter $q.HardcodedFilters .Name)}}
	{{camel .Name}} :=  {{ getHardcodedValue $q.HardcodedFilters .Name .Type}}
		{{- end}}
	{{- end}}

	filters := &{{.ReturnType.Module.Store -}}.{{- title .ReturnType.Name -}}Filters{
		{{- template "getargs" . }}
	}
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err :=  store.{{- title .ReturnType.Module.Name -}}Store.Get{{- plural (title .ReturnType.Name)}} (
		ctx,
		filters,
		in.After,
		in.Before,
		first,
		last, 
		in.SortBy,
		in.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	return New{{- .ReturnType.ConnectionName}}(filters, data, hasNextPage, hasPreviousPage, &firstCursor, &lastCursor), nil
}
{{- else}} 
{{- if .Query.Parent }}
func (n *{{- title .Query.Parent.Name -}}) {{title .Query.Name}}(ctx context.Context {{- if (ne (len .Query.Args) 0)}}, in *{{title .Query.Name -}}Input {{- end}}) (*{{- title .ReturnType.Name}}, error) {
{{- else }}
func (r *Resolver) {{title .Query.Name}}(ctx context.Context {{- if (ne (len .Query.Args) 0)}}, in *{{title .Query.Name -}}Input {{- end}}) (*{{- title .ReturnType.Name }}, error) {
{{- end}}
	
	filters := &{{.ReturnType.Module.Store -}}.{{- title .ReturnType.Name -}}Filters{
		{{- template "getargs" . }}
	}
	store := storemanager.FromContext(ctx)
	{{.ReturnType.Name }}, err := store.{{- title .ReturnType.Module.Name -}}Store.GetOne{{- title .ReturnType.Name}}(ctx, filters)
	if err != nil {
		return nil, err
	}
	return New{{- title .ReturnType.Name }}({{.ReturnType.Name }}), nil
}
{{- end}}

{{- define "getargs"}}
{{- $q := .Query -}}
{{- if (isElemOfStrArray ($q.ArgNames) "ids")}}
	Ids: in.Ids,
{{- end}}
{{- $t := .}}
{{- range $q.Args}}
	{{- if (isElemOfStrArray ($t.ReturnType.Filters) .Name)}}
		{{- if (isHardCodedFilter $q.HardcodedFilters .Name)}}
		{{title .Name}} : &{{- camel .Name}},
		{{- else}}
		{{title .Name}} :in.{{- title .Name}},
		{{- end}}
	{{- end}}
{{- end}}
{{- range .ReturnType.ReferedTypes}}
	{{- if (and $q.Parent (eq .Name $q.Parent.Name))}}
	{{title .Name}}: n.ref,
	{{- end}}
{{- end}}
{{- end}}

`
