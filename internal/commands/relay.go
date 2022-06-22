package commands

import (
	"context"
	"flag"
	"io/fs"
	"os"

	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/generators/relayjs"
	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/parser"
	log "github.com/sirupsen/logrus"
)

// Relay is the command doing code generation
type Relay struct {
	*CommandWriter
	flags    flag.FlagSet
	filepath string
	verbose  bool
	help     bool
	fs.FS
	outputPath string
}

// NewRelay return a new instance of Relay
func NewRelay(writer *CommandWriter) *Relay {
	return &Relay{
		CommandWriter: writer,
		FS:            os.DirFS("."),
		outputPath:    ".",
	}
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
func (r *Relay) Help() {
	helpText := `
	matro Relay  f [--file] <file path>
	It generates all Relay side code
	[vv] for verbose

	`
	r.Write(helpText)
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

	matroConfigs, err := c.FS.Open(c.filepath)
	if err != nil {
		return err
	}
	d, err := parser.NewDefinition().ParseFrom(matroConfigs)
	if err != nil {
		return err
	}
	path := "."
	typeDefs, err := types.GetTypeDefs(d)
	if err != nil {
		return err
	}
	generators := []generators.IGenerator{
		relayjs.NewGenerator(path, d, typeDefs),
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
