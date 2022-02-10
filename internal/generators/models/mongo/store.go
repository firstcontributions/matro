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
