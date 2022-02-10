package mongo

var storeInterfaceTpl = `
package {{ .Name -}}store



type Store interface {
	{{- range .Types}}
	{{- if .IsNode }}
	// {{ .Name }} methods
	Create{{- title .Name -}} (ctx context.Context, {{ .Name}}  *{{- title .Name}}) (*{{- title .Name}}, error)
	Get{{- title .Name -}}ByID (ctx context.Context, id string) (*{{- title .Name}}, error)
	Get{{- title (plural .Name) -}} (ctx context.Context, ids []string,
		{{- if not (empty .SearchFields) -}}
		search *string,
		{{- end -}}
		{{- template "getargs" . -}}
		after *string, before *string, first *int64, last *int64) ([]* {{- title .Name}}, bool, bool,string, string, error) 

	Update{{- title .Name -}} (ctx context.Context, id string, update *{{- title .Name -}}Update) (error) 
	Delete{{- title .Name -}}ByID (ctx context.Context, id string) (error)
	{{- end}}
	{{- end}}
}


{{- define "getargs" -}}
{{- $t := . -}}
{{- range .Filters -}}
	{{ . }} *{{$t.FieldType . -}},
{{- end -}}
{{- range .ReferedFields -}}
	{{.}} *string,
{{- end -}}
{{- end -}}
`
