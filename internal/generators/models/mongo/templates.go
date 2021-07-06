package mongo

const storeTpl = `
package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DB{{title (plural .Name) -}} = {{- plural .Name }}
	{{- range .Types }}
	Collection{{title (plural .Name)}} = {{- plural .Name }}
	{{- end}}
)

type {{ title (plural .Name) -}}Store struct {
	client *mongo.Client
}

// New {{- title (plural .Name) -}}Store makes connection to mongo server by provided url 
// and return an instance of the client
func New {{- title (plural .Name) -}}Store(ctx context.Context, mongoUrl string) (* {{ title (plural .Name) -}}Store, error) {
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
	return &{{- title (plural .Name) -}}Store {
		client: client,
	}, nil
} 

func (s *{{- title (plural .Name) -}}Store) getCollection (collection string) *mongo.Collection {
	return s.client.Database(DB{{ title (plural .Name) -}}).Collection(collection)
}
`

const modelTpl = `
package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"

)

func (s *{{- title (plural .Module) -}}Store) Get{{- title .Name -}}ByID (ctx types.Context, id string) (* {{plural .Module -}}store. {{- title .Name}}, error) {
	query := bson.M{
		"_id": id,
	}
	var {{.Name}} {{plural .Module -}}store. {{- title .Name}}
	if err := s.getCollection(Collection{{title (plural .Name)}}).FindOne(ctx, query).Decode(&{{- .Name -}}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &{{- .Name}}, nil
}

func (s *{{- title (plural .Module) -}}Store) Update{{- title .Name -}} (ctx types.Context, {{.Name}} *{{- plural .Module -}}store. {{- title .Name}}) (* {{plural .Module -}}store. {{- title .Name}}, error) {
	query := bson.M{
		"_id": {{.Name -}}.Id,
	}
	var {{.Name}} {{ title .Name}}
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
package {{plural .Module -}}store

type {{title .Name}} struct {
	{{- counter 0}} 
	{{- range .Fields}}
	{{- if  (not (and .IsJoinedData  .IsList))}}
	{{ .GoName}}  {{- .GoType}}` + "`bson:\"{{- .Name}}\"`" + `  
	{{- end}}
	{{- end}}
}`
