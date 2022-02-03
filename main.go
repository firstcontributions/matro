package main

import (
	"context"
	"log"

	"github.com/firstcontributions/matro/internal/ctxkeys"
	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/firstcontributions/matro/pkg/spinner"
)

func main() {
	d, err := parser.NewDefinition().ParseFromFile("./input.json")
	if err != nil {
		log.Fatal(err)
	}
	generatorTypes := []generators.Type{
		generators.TypeGQLSchema,
		generators.TypeGQLServer,
		generators.TypeGRPCProto,
		generators.TypeGRPCStore,
		generators.TypeMongoModel,
		generators.TypeGoMod,
	}

	if err := generate(d, generatorTypes); err != nil {
		// this can be a pretty print for the final version
		log.Fatal(err)
	}
}

// generate generates code for give types refered by parsed format of input json
func generate(d *parser.Definition, generatorTypes []generators.Type) error {
	// this context can be used to close the spinner
	ctx, cancel := context.WithCancel(context.Background())
	// make sure context is cancelled once the code generation completed
	defer cancel()
	s := spinner.NewSpinner("generating")
	go s.Start(ctx)

	// this needs to be taken as an command line argument, hardcoding for now
	// the default value can be $(pwd)
	// it wont be handy to use $(pwd) as default in development time
	basePath := "./"
	ctx = context.WithValue(ctx, ctxkeys.Spinner, s)
	for _, gt := range generatorTypes {
		g := generators.GetGenerator(gt, basePath, d)
		// will terminate all generations if any of the generators are
		// throwing an error
		if err := g.Generate(ctx); err != nil {
			return err
		}
	}
	return nil
}
