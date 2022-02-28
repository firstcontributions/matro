package main

import (
	"context"
	"log"

	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/generators/gomod"
	"github.com/firstcontributions/matro/internal/generators/graphql/gocode"
	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/generators/grpc/proto"
	"github.com/firstcontributions/matro/internal/generators/grpc/store"
	"github.com/firstcontributions/matro/internal/generators/models/mongo"
	"github.com/firstcontributions/matro/internal/generators/relayjs"
	"github.com/firstcontributions/matro/internal/parser"
)

func main() {
	d, err := parser.NewDefinition().ParseFromFile("./input.json")
	if err != nil {
		log.Fatal(err)
	}
	path := "."
	generators := []generators.IGenerator{
		schema.NewGenerator(path, d),
		gocode.NewGenerator(path, d),
		proto.NewGenerator(path, d),
		store.NewGenerator(path, d),
		mongo.NewGenerator(path, d),
		relayjs.NewGenerator(path, d),
		gomod.NewGenerator(path, d),
		relayjs.NewGenerator(path, d),
	}
	ctx := context.Background()
	for _, g := range generators {
		// will terminate all generations if any of the generators are
		// throwing an error
		if err := g.Generate(ctx); err != nil {
			log.Fatal(err)
		}
	}
}
