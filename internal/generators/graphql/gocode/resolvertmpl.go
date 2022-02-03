package gocode

var resolverTmpl = `
package schema

import (
	{{- $g := .}}
	{{- range .Modules }}
	"{{- $g.Repo -}}/internal/models/{{- .Name -}}store"
	graphql "github.com/graph-gophers/graphql-go"
	{{- end}}
)

type Store struct {
	{{- range .Modules }} 
	{{ .Name -}}Store {{ .Name -}}store.Store
	{{- end}}
}

func NewStore(
	{{- range .Modules }} 
	{{ .Name -}}Store {{ .Name -}}store.Store,
	{{- end}}
) *Store {
	return &Store{
		{{- range .Modules }} 
		{{ .Name -}}Store :{{ .Name -}}Store,
		{{- end}}
	}
}

type Resolver struct {
	*Store
}


func (r *Resolver) Viewer(ctx context.Context) (*User, error) {
	id := ctx.Value("user_id").(string)

	data, err := r.usersStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewUser(data), nil
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
	HasNextPage bool
	HasPreviousPage bool
	StartCursor *string
	EndCursor *string
}

func ParseGraphqlID(gid graphql.ID) (*IDMarshaller, error) {
	sDec, err := base64.StdEncoding.DecodeString(string(gid))
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

func (id *IDMarshaller) String() string {
	return base64.StdEncoding.EncodeToString(
		[]byte(id.Type + ":" + id.ID),
	)
}

func (id *IDMarshaller) ToGraphqlID() graphql.ID {
	return graphql.ID(id.String())
}
`
