package gocode

const mutationTmpl = `
package schema


func (m *Resolver) Create{{- title .Name }}(
	ctx context.Context, 
	args struct {
		{{title .Name}}  *Create{{- title .Name -}}Input
	},
) (*{{- title .Name }}, error) {
	store := storemanager.FromContext(ctx)

	{{- $type := .}}
	{{.Name}}, err := store.{{- title .Module.Name -}}Store.Create{{- title .Name }}(ctx, args.{{title .Name -}}.ToModel())
	if err != nil {
		return nil, err
	}
	return New{{- title .Name }}({{.Name}}), nil
}
`
