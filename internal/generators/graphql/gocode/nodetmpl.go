package gocode

var nodeTmpl = `
package schema

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"
)

type Node interface {
	ID(context.Context) graphql.ID
}

type NodeResolver struct {
	Node
}

type NodeIDInput struct {
	ID graphql.ID
}

func (r *Resolver) Node(ctx context.Context, in NodeIDInput) (*NodeResolver, error) {
	id, err := ParseGraphqlID(in.ID)
	if err != nil {
		return nil, err
	}
	switch id.Type {
		{{- range .Types}}
		{{- if .IsNode }}
	case "{{- .Name -}}":
		{{.Name}}Data, err := r.{{- .Module -}}Store.Get{{- title .Name -}}ByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		{{.Name -}}Node := New{{- title .Name}}({{.Name}}Data)
		return &NodeResolver{
			Node: {{.Name -}}Node,
		}, nil
		{{- end}}
		{{- end}}
	}
	return nil, errors.New("invalid ID")
}

{{- range .Types}}
{{- if .IsNode }}
func (n *NodeResolver) To{{-  title .Name -}}() (*{{-  title .Name -}}, bool) {
	t, ok := n.Node.(*{{-  title .Name -}})
	return t, ok
}
{{- end}}
{{- end}}


`
