package mongo

const storeManagerTmpl = `
package storemanager

import (
	{{- $g := .}}
	{{- range .Modules }}
	"{{- $g.Repo -}}/internal/models/{{- .Name -}}store"
	{{- end}}
)


type storeCtxKey int

const (
	store storeCtxKey = iota
)

type Store struct {
	{{- range .Modules }} 
	{{ title .Name -}}Store {{ .Name -}}store.Store
	{{- end}}
}

func NewStore(
	{{- range .Modules }} 
	{{ .Name -}}Store {{ .Name -}}store.Store,
	{{- end}}
) *Store {
	return &Store{
		{{- range .Modules }} 
		{{ title .Name -}}Store :{{ .Name -}}Store,
		{{- end}}
	}
}


func ContextWithStore(ctx context.Context, s *Store) context.Context {
	return context.WithValue(ctx, store, s)
}

func FromContext(ctx context.Context) *Store {
	return ctx.Value(store).(*Store)
}

`
