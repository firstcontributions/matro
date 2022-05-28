package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/generators/relayjs"
	"github.com/firstcontributions/matro/internal/parser"
	log "github.com/sirupsen/logrus"
)

// Relay is the command doing code generation
type Relay struct {
	flags    flag.FlagSet
	filepath string
	verbose  bool
	help     bool
}

// NewRelay return a new instance of Relay
func NewRelay() *Relay {
	return &Relay{}
}

// InitFlags will initialize all flags
func (c *Relay) InitFlags() {
	c.flags.StringVar(&c.filepath, "f", "matro.json", "file path")
	c.flags.StringVar(&c.filepath, "file", "matro.json", "file path")
	c.flags.BoolVar(&c.help, "h", false, "help")
	c.flags.BoolVar(&c.help, "help", false, "help")
	c.flags.BoolVar(&c.verbose, "vv", false, "verbose")

}

// ParseFlags will parse given flags
func (c *Relay) ParseFlags(args []string) {
	c.flags.Parse(args)
}

// Help prints the help message
func (Relay) Help() {
	helpText := `
	matro Relay  f [--file] <file path>
	It generates all Relay side code
	[vv] for verbose
	`
	fmt.Println(helpText)
}

// Exec will execute the core command functionality, here it Relays and saves the code
func (c *Relay) Exec() error {
	if c.help {
		c.Help()
		return nil
	}
	if c.verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.FatalLevel)
	}
	d, err := parser.NewDefinition().ParseFromFile("./input.json")
	if err != nil {
		return err
	}
	path := "."
	generators := []generators.IGenerator{
		relayjs.NewGenerator(path, d),
	}
	ctx := context.Background()
	for _, g := range generators {
		// will terminate all generations if any of the generators are
		// throwing an error
		if err := g.Generate(ctx); err != nil {
			return err
		}
	}
	return nil
}
