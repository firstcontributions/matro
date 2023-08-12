package gocode

const mutationTmpl = `
package schema


func (m *Resolver) Create{{- title .Name }}(
	ctx context.Context, 
	args struct {
		{{title .Name}}  *Create{{- title .Name -}}Input
	},
) (*{{- title .Name }}, error) {	
	session := session.FromContext(ctx)
	if session == nil {
		return nil, errors.New("unauthorized")
	}

	{{ $type := .}}
	{{.Name}}ModelInput, err := args.{{title .Name -}}.ToModel()
	if err != nil {
		return nil, err
	}
	{{- range .Fields}}
	{{- if .ViewerRefence}} 
	{{$type.Name -}}ModelInput.{{- title .Name}} = session.UserID()
	{{- end}}
	{{- end}}

	{{- range .ReferedTypes}}
	{{- if .IsViewerType}} 
	{{$type.Name -}}ModelInput.{{- title .Name}}ID = session.UserID()
	{{- end}}
	{{- end}}



	ownership := &authorizer.Scope{
		Users: []string{session.UserID()},
	}
	{{.Name}}, err := storemanager.FromContext(ctx).
		{{- title .Module.Name -}}Store.Create{{- title .Name }}(ctx, {{.Name}}ModelInput, ownership)
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
	session := session.FromContext(ctx)
	if session == nil {
		return nil, errors.New("unauthorized")
	}

	store := storemanager.FromContext(ctx)

	id, err := ParseGraphqlID(args.{{title .Name -}}.ID)
	if err != nil {
		return nil, err
	}

	{{.Name}}, err := store.{{- title .Module.Name -}}Store.Get{{- title .Name }}ByID(ctx, id.ID)
	if err != nil {
		return nil, err
	}

	if !authorizer.IsAuthorized(session.Permissions, {{.Name -}}.Ownership, authorizer.{{title .Name}}, authorizer.OperationUpdate) {
		return nil, errors.New("forbidden")
	}
	{{- $type := .}}
	if err := store.{{- title .Module.Name -}}Store.Update{{- title .Name }}(ctx, id.ID, args.{{title .Name -}}.ToModel());err != nil {
		return nil, err
	}
	{{.Name}}, err = store.{{- title .Module.Name -}}Store.Get{{- title .Name }}ByID(ctx, id.ID)
	if err != nil {
		return nil, err
	}
	return New{{- title .Name }}({{.Name}}), nil
}
{{- end}}
`
