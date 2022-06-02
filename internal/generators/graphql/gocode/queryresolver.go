package gocode

var queryResolverTmpl = `
package schema

{{- if (ne (len .Query.Args) 0)}}
type {{title .Query.Name -}}Input struct {
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
func (n *{{- title .Query.Parent.Name -}}) {{title .Query.Name}}(ctx context.Context, in *{{title .Query.Name -}}Input) (*{{.ReturnType.ConnectionName}}, error) {
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
	store := storemanager.FromContext(ctx)
	{{- $q := .Query}}
	{{- range .Query.Args}}
		{{- if (isHardCodedFilter $q.HardcodedFilters .Name)}}
	{{camel .Name}} :=  {{ getHardcodedValue $q.HardcodedFilters .Name .Type}}
		{{- end}}
	{{- end}}
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err :=  store.{{- title .ReturnType.Module.Name -}}Store.Get{{- plural (title .ReturnType.Name)}} (
		ctx,
		{{- template "getargs" .}}
	)
	if err != nil {
		return nil, err
	}
	return New{{- .ReturnType.ConnectionName}}(data, hasNextPage, hasPreviousPage, &firstCursor, &lastCursor), nil
}
{{- else}} 
{{- if .Query.Parent }}
func (n *{{- title .Query.Parent.Name -}}) {{title .Query.Name}}(ctx context.Context{{- if (ne (len .Query.Args) 0)}}, in *{{title .Query.Name -}}Input {{- end}}) (*{{- title .ReturnType.Name}}, error) {
{{- else }}
func (r *Resolver) {{title .Query.Name}}(ctx context.Context, {{- if (ne (len .Query.Args) 0)}}, in *{{title .Query.Name -}}Input {{- end}}) (*{{- title .ReturnType.Name }}, error) {
{{- end}}
	return New{{- title .ReturnType.Name }}(), nil
}
{{- end}}

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
{{- $t := .}}
{{- range $q.Args}}
	{{- if (isElemOfStrArray ($t.ReturnType.Filters) .Name)}}
		{{- if (isHardCodedFilter $q.HardcodedFilters .Name)}}
		&{{- camel .Name}},
		{{- else}}
		in.{{- title .Name}},
		{{- end}}
	{{- end}}
{{- end}}
	n.ref,
	in.After,
	in.Before,
	first,
	last, 
{{- end}}

`
