package commands

import (
	"context"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/parser"

	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/generators/gomod"
	"github.com/firstcontributions/matro/internal/generators/graphql/gocode"
	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/generators/grpc/proto"
	"github.com/firstcontributions/matro/internal/generators/grpc/service"
	"github.com/firstcontributions/matro/internal/generators/grpc/store"
	"github.com/firstcontributions/matro/internal/generators/models/mongo"
)

// Server is the command doing code generation
type Server struct {
	*CodeGenerator
}

// NewServer return a new instance of Server
func NewServer(writer *CommandWriter) *Server {
	return &Server{
		NewCodeGenerator(writer),
	}
}

// Help prints the help message
func (c *Server) Help() {
	helpText := `
	matro server  -f [--file] <file path>
	It generates all server side code
	[-vv] for verbose

	`
	c.Write(helpText)
}

// get genertors will return an instance of each server code generator
func (c *Server) getGenerators(d *parser.Definition, typeDefs *types.TypeDefs) []generators.IGenerator {
	return []generators.IGenerator{
		gomod.NewGenerator(c.outputPath, d),
		schema.NewGenerator(c.outputPath, d, typeDefs),
		gocode.NewGenerator(c.outputPath, d, typeDefs),
		proto.NewGenerator(c.outputPath, d, typeDefs),
		store.NewGenerator(c.outputPath, d, typeDefs),
		mongo.NewGenerator(c.outputPath, d, typeDefs),
		service.NewGenerator(c.outputPath, d, typeDefs),
	}
}

// Exec will execute the code generation for all given generators based on given configs
func (c *Server) Exec() error {
	if countinue := c.Setup(); !countinue {
		return nil
	}
	d, typeDefs, err := c.GetDefenitionsAndTypes()
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, g := range c.getGenerators(d, typeDefs) {
		// will terminate all generations if any of the generators are
		// throwing an error
		if err := g.Generate(ctx); err != nil {
			return err
		}
	}
	return nil
}
