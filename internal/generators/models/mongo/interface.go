package mongo

var storeInterfaceTpl = `
package {{ .Name -}}store



type Store interface {
	{{- range .Types}}
	{{- if .IsNode }}
	// {{ .Name }} methods
	Get{{- title .Name -}}ByID (ctx context.Context, id string) (*{{- title .Name}}, error)
	GetOne{{- title .Name -}} (ctx context.Context, filters *{{- title .Name -}}Filters) (*{{- title .Name}}, error)
	Get{{- title (plural .Name) -}} (ctx context.Context, filters *{{- title .Name -}}Filters, after *string, before *string, first *int64, last *int64, sortBy {{title .Name -}}SortBy, sortOrder *string) ([]* {{- title .Name}}, bool, bool, []string, error) 
	Count{{- title (plural .Name) -}} (ctx context.Context, filters *{{- title .Name -}}Filters) (int64, error) 
	{{- if (ne .Module.DB "") }}
	Create{{- title .Name -}} (ctx context.Context, {{ .Name}}  *{{- title .Name}}) (*{{- title .Name}}, error)
	Update{{- title .Name -}} (ctx context.Context, id string, update *{{- title .Name -}}Update) (error) 
	Delete{{- title .Name -}}ByID (ctx context.Context, id string) (error)
	{{- end}}
	{{- end}}
	{{- end}}
}


{{- define "getargs" -}}
{{- $t := . -}}
{{- range .Filters -}}
	{{ . }} *{{$t.FieldType . -}},
{{- end -}}
{{- range .ReferedTypes -}}
	{{- if (eq $t.Module.Name .Module.Name)}}
	{{.Name}} *{{- title .Name}},
	{{- else}}
	{{.Name}} *{{.Module.Store}}.{{- title .Name}},
	{{- end}}
{{- end -}}
{{- end -}}
`
