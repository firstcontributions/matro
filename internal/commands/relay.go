package commands

import (
	"context"

	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/generators/relayjs"
	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/parser"
)

// Relay is the command doing code generation
type Relay struct {
	*CodeGenerator
}

// NewRelay return a new instance of Relay
func NewRelay(writer *CommandWriter) *Relay {
	return &Relay{
		NewCodeGenerator(writer),
	}
}

// Help prints the help message
func (c *Relay) Help() {
	helpText := `
	matro Relay  f [--file] <file path>
	It generates all Relay side code
	[vv] for verbose

	`
	c.Write(helpText)
}

// get genertors will return an instance of each server code generator
func (c *Relay) getGenerators(d *parser.Definition, typeDefs *types.TypeDefs) []generators.IGenerator {
	return []generators.IGenerator{
		relayjs.NewGenerator(c.outputPath, d, typeDefs),
	}
}

// Exec will execute the core command functionality, here it Relays and saves the code
func (c *Relay) Exec() error {
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
