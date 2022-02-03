package mongo

var storeInterfaceTpl = `
package {{ .Name -}}store



type Store interface {
	{{- range .Types}}
	{{- if .IsNode }}
	// {{ .Name }} methods
	Create{{- title .Name -}} (context.Context,  *{{- title .Name}}) (*{{- title .Name}}, error)
	Get{{- title .Name -}}ByID (context.Context, string) (*{{- title .Name}}, error)
	Get{{- title (plural .Name) -}} (context.Context, []string,
		{{- if not (empty .SearchFields) -}}
		*string,
		{{- end -}}
		{{- template "getargs" . -}}
		*string, *string,  *int64,  *int64) ([]* {{- title .Name}}, bool, bool,string, string, error) 

	{{- if .Mutatable}}
	Update{{- title .Name -}} (context.Context, string, *{{- title .Name -}}Update) (error) 
	{{- end}}
	Delete{{- title .Name -}}ByID (context.Context, string) (error)
	{{- end}}
	{{- end}}
}


{{- define "getargs" -}}
{{- $t := . -}}
{{- range .Filters -}}
	*{{$t.FieldType . -}},
{{- end -}}
{{- range .ReferedFields -}}
	*string,
{{- end -}}
{{- end -}}
`
