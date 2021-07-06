package main

import (
	"context"

	"github.com/firstcontributions/matro/internal/cleaner"
	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/firstcontributions/matro/pkg/spinner"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := spinner.NewSpinner(ctx, "generating")
	go s.Start()
	d, err := parser.NewDefinition().ParseFromFile("./input.json")
	if err != nil {
		panic(err)
	}
	basePath := "./__generated"
	s.Update("cleaning code")
	cleaner.Clean(basePath)
	s.Update("graphql schema")
	gs := generators.GetGenerator(basePath, "schema", d)
	err = gs.Generate()
	if err != nil {
		panic(err)
	}
	s.Update("graphql server")
	gc := generators.GetGenerator(basePath, "gocode", d)
	err = gc.Generate()
	if err != nil {
		panic(err)
	}
	s.Update("grpc protobuf")
	gp := generators.GetGenerator(basePath, "proto", d)
	err = gp.Generate()
	if err != nil {
		panic(err)
	}

	s.Update("mongo models")
	gm := generators.GetGenerator(basePath, "mongo", d)
	err = gm.Generate()
	if err != nil {
		panic(err)
	}
}
