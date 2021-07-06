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

const modelTpl = `
package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"{{- .Repo -}}/internal/models/{{- .Module -}}store"


)

func (s *{{- title .Module -}}Store) Get{{- title .Name -}}ByID (ctx context.Context, id string) (* {{ .Module -}}store. {{- title .Name}}, error) {
	query := bson.M{
		"_id": id,
	}
	var {{.Name}} {{ .Module -}}store. {{- title .Name}}
	if err := s.getCollection(Collection{{title (plural .Name)}}).FindOne(ctx, query).Decode(&{{- .Name -}}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &{{- .Name}}, nil
}

func (s *{{- title .Module -}}Store) Update{{- title .Name -}} (ctx context.Context, {{.Name}} *{{-  .Module -}}store. {{- title .Name}}) (* {{ .Module -}}store. {{- title .Name}}, error) {
	query := bson.M{
		"_id": {{.Name -}}.Id,
	}
	if _, err := s.getCollection(Collection{{title (plural .Name)}}).UpdateOne(ctx, query, {{.Name}}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return {{ .Name}}, nil
}
`

const modelTyp = `
package {{ .Module -}}store

type {{title .Name}} struct {
	{{- counter 0}} 
	{{- range .Fields}}
	{{- if  (not (and .IsJoinedData  .IsList))}}
	{{ .GoName}}  {{- .GoType}}` + "`bson:\"{{- .Name}}\"`" + `  
	{{- end}}
	{{- end}}
}`
