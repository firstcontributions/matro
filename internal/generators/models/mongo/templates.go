package mongo

const storeTpl = `
package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DB{{title (plural .Name) -}} = "{{- plural .Name }}"
	{{- range .Types }}
	Collection{{title (plural .Name)}} = "{{- plural .Name }}"
	{{- end}}
)

type {{ title .Name -}}Store struct {
	client *mongo.Client
}

// New {{- title .Name -}}Store makes connection to mongo server by provided url 
// and return an instance of the client
func New {{- title .Name -}}Store(ctx context.Context, url string) (* {{ title .Name -}}Store, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err 
	}
	return &{{- title .Name -}}Store {
		client: client,
	}, nil
} 

func (s *{{- title .Name -}}Store) getCollection (collection string) *mongo.Collection {
	return s.client.Database(DB{{ title (plural .Name) -}}).Collection(collection)
}
`

const crudTpl = `
package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gokultp/go-mongoqb"
	"{{- .Repo -}}/internal/models/{{- .Module -}}store"


)

func (s *{{- title .Module -}}Store) Create{{- title .Name -}} (ctx context.Context, {{.Name}} *{{-  .Module -}}store. {{- title .Name}}) (* {{ .Module -}}store. {{- title .Name}}, error) {
	now := time.Now()
	{{.Name -}}.TimeCreated = &now
	{{.Name -}}.TimeUpdated = &now
	if _, err := s.getCollection(Collection{{title (plural .Name)}}).InsertOne(ctx, {{.Name}}); err != nil {
		return nil, err
	}
	return {{ .Name}}, nil
}

func (s *{{- title .Module -}}Store) Get{{- title .Name -}}ByID (ctx context.Context, id string) (* {{ .Module -}}store. {{- title .Name}}, error) {
	qb := mongoqb.NewQueryBuilder().
			Eq("_id", id)
	var {{.Name}} {{ .Module -}}store. {{- title .Name}}
	if err := s.getCollection(Collection{{title (plural .Name)}}).FindOne(ctx, qb.Build()).Decode(&{{- .Name -}}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &{{- .Name}}, nil
}

func (s *{{- title .Module -}}Store) Get{{- title (plural .Name) -}} (
	ctx context.Context,
	{{- if not (empty .SearchFields)}}
	search *string,
	{{- end}}
	{{- template "getargs" . }}
	offset *string,
	limit *int, 
) (
	[]*{{ .Module -}}store. {{- title .Name}}, 
	error,
) {
	qb := mongoqb.NewQueryBuilder()
	{{- range .Filters }}
	qb.Eq("{{- . -}}", {{.}})
	{{- end }}

	var {{plural .Name}} []*{{ .Module -}}store. {{- title .Name}}
	cursor, err := s.getCollection(Collection{{title (plural .Name)}}).Find(ctx, qb.Build())
	if  err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &{{- plural .Name}})
	if err != nil {
		return nil, err
	}
	return {{ plural .Name}}, nil
}
{{- if .Mutatable}}
func (s *{{- title .Module -}}Store) Update{{- title .Name -}} (ctx context.Context, id string, {{.Name -}}Update *{{-  .Module -}}store. {{- title .Name -}}Update) (error) {
	qb := mongoqb.NewQueryBuilder().
			Eq("_id", id)

	now := time.Now()
	{{.Name -}}Update.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
	SetFields({{- .Name -}}Update)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(Collection{{title (plural .Name)}}).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}
{{- end}}
func (s *{{- title .Module -}}Store) Delete{{- title .Name -}}ByID (ctx context.Context, id string) (error) {
	qb := mongoqb.NewQueryBuilder().
			Eq("_id", id)
	if _,  err := s.getCollection(Collection{{title (plural .Name)}}).DeleteOne(ctx, qb.Build()); err != nil {
		return  err
	}
	return nil
}

{{- define "getargs"}}
{{- $t := .}}
{{- range .Filters}}
	{{.}} *{{$t.FieldType .}},
{{- end}}
{{- end}}
`

const modelTyp = `
package {{ .Module -}}store

type {{title .Name}} struct {
	{{- counter 0}} 
	{{- range .Fields}}
	{{- if  (not (and .IsJoinedData  .IsList))}}
	{{ .GoName true}}  {{- .GoType }}` + "`bson:\"{{- .Name}}\"`" + `  
	{{- end}}
	{{- end}}
}

{{- if .Mutatable}}
type {{title .Name -}}Update struct {
	{{- counter 0}} 
	{{- range .Fields}}
	{{- if  .IsMutatable }}
	{{ .GoName true}}  {{- .GoType }}` + "`bson:\"{{- .Name}}\"`" + `  
	{{- end}}
	{{- end}}
}
{{- end}}
`

var storeInterfaceTpl = `
package {{ .Name -}}store



type Store interface {
	{{- range .Types}}

	// {{ .Name }} methods
	Create{{- title .Name -}} (context.Context,  *{{- title .Name}}) (*{{- title .Name}}, error)
	Get{{- title .Name -}}ByID (context.Context, string) (*{{- title .Name}}, error)
	Get{{- title (plural .Name) -}} (context.Context,
		{{- if not (empty .SearchFields) -}}
		*string,
		{{- end -}}
		{{- template "getargs" . -}}
		*string,*int) ([]* {{- title .Name}}, error) 

	{{- if .Mutatable}}
	Update{{- title .Name -}} (context.Context, *{{- title .Name -}}Update) (error) 
	{{- end}}
	Delete{{- title .Name -}}ByID (context.Context, string) (error)
	{{- end}}
}


{{- define "getargs" -}}
{{- $t := . -}}
{{- range .Filters -}}
	*{{$t.FieldType . -}},
{{- end -}}
{{- end -}}
`
