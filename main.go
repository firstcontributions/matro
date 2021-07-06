package main

import (
	"context"

	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/firstcontributions/matro/pkg/spinner"
)

func main() {
	d, err := parser.NewDefinition().ParseFromFile("./input.json")
	if err != nil {
		panic(err)
	}
	generatorTypes := []generators.Type{
		generators.TypeGQLSchema,
		generators.TypeGQLServer,
		generators.TypeGRPCProto,
		generators.TypeMongoModel,
	}

	if err := generate(d, generatorTypes); err != nil {
		panic(err)
	}
}

func generate(d *parser.Definition, generatorTypes []generators.Type) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := spinner.NewSpinner(ctx, "generating")
	go s.Start()
	basePath := "./__generated"

	for _, gt := range generatorTypes {
		s.Update(string(gt))
		g := generators.GetGenerator(gt, basePath, d)
		if err := g.Generate(); err != nil {
			return err
		}
	}
	return nil
}
