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

{{- if  .GraphqlOps.Update }}
func (m *Resolver) Update{{- title .Name }}(
	ctx context.Context, 
	args struct {
		{{title .Name}}  *Update{{- title .Name -}}Input
	},
) (*{{- title .Name }}, error) {
	store := storemanager.FromContext(ctx)

	id, err := ParseGraphqlID(args.{{title .Name -}}.ID)
	if err != nil {
		return nil, err
	}
	{{- $type := .}}
	if err := store.{{- title .Module.Name -}}Store.Update{{- title .Name }}(ctx, id.ID, args.{{title .Name -}}.ToModel());err != nil {
		return nil, err
	}
	user, err := store.{{- title .Module.Name -}}Store.Get{{- title .Name }}ByID(ctx, id.ID)
	if err != nil {
		return nil, err
	}
	return New{{- title .Name }}({{.Name}}), nil
}
{{- end}}
`
