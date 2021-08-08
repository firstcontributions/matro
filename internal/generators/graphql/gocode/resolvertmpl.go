package gocode

var resolverTmpl = `
package schema

import (
	{{- $g := .}}
	{{- range .Modules }}
	"{{- $g.Repo -}}/internal/models/{{- .Name -}}store"
	{{- end}}
)

type Resolver struct {
	{{- range .Modules }} 
	{{ .Name -}}Store {{ .Name -}}store.Store
	{{- end}}
}


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

type PageInfo struct {
	HasNextPage *bool
	HasPreviousPage *bool
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
