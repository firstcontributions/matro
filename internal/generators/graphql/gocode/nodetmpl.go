package gocode

var nodeTmpl = `
package schema

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"
)

type Node interface {
	ID(context.Context) *graphql.ID
}

type NodeResolver struct {
	Node
}

type NodeIDInput struct {
	ID *graphql.ID
}

{{- range .Types}}
{{- if .IsNode }}
func (n *NodeResolver) To{{-  title .Name -}}() (*{{-  title .Name -}}, bool) {
	t, ok := n.Node.(*{{-  title .Name -}})
	return t, ok
}
{{- end}}
{{- end}}


type IDMarshaller struct {
	Type string
	ID string
}

func NewIDMarshaller(t, id string) *IDMarshaller {
	return &IDMarshaller{
		Type:t,
		ID: id,
	}
}

func ParseGraphqlID(gid *graphql.ID) (*IDMarshaller, error) {
	if gid == nil {
		return nil, errors.New("empty ID")
	}
	sDec, err := base64.StdEncoding.DecodeString(string(*gid))
	if err != nil {
		return nil, errors.New("invalid ID")
	}
	ids := strings.Split(string(sDec), ":")
	if len(ids) != 2 {
		return nil, errors.New("invalid ID")
	}
	return &IDMarshaller {
		Type: ids[0],
		ID: ids[1],
	}, nil
}

func (id *IDMarshaller) ToGraphqlID() *graphql.ID {
	encoded := base64.StdEncoding.EncodeToString(
		[]byte(id.Type + ":" + id.ID),
	)
	gid := graphql.ID(encoded)
	return &gid
}
`
