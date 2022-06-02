package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/generators/gomod"
	"github.com/firstcontributions/matro/internal/generators/graphql/gocode"
	"github.com/firstcontributions/matro/internal/generators/graphql/schema"
	"github.com/firstcontributions/matro/internal/generators/grpc/proto"
	"github.com/firstcontributions/matro/internal/generators/grpc/service"
	"github.com/firstcontributions/matro/internal/generators/grpc/store"
	"github.com/firstcontributions/matro/internal/generators/models/mongo"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/sirupsen/logrus"
)

// Server is the command doing code generation
type Server struct {
	flags    flag.FlagSet
	filepath string
	verbose  bool
	help     bool
}

// NewServer return a new instance of Server
func NewServer() *Server {
	return &Server{}
}

// InitFlags will initialize all flags
func (c *Server) InitFlags() {
	c.flags.StringVar(&c.filepath, "f", "matro.json", "file path")
	c.flags.StringVar(&c.filepath, "file", "matro.json", "file path")
	c.flags.BoolVar(&c.help, "h", false, "help")
	c.flags.BoolVar(&c.help, "help", false, "help")
	c.flags.BoolVar(&c.verbose, "vv", false, "verbose")
}

// ParseFlags will parse given flags
func (c *Server) ParseFlags(args []string) {
	c.flags.Parse(args)
}

// Help prints the help message
func (Server) Help() {
	helpText := `
	matro server  -f [--file] <file path>
	It generates all server side code
	[-vv] for verbose
	`
	fmt.Println(helpText)
}

// Exec will execute the core command functionality, here it Servers and saves the code
func (c *Server) Exec() error {
	if c.help {
		c.Help()
		return nil
	}
	if c.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.FatalLevel)
	}
	d, err := parser.NewDefinition().ParseFromFile(c.filepath)
	if err != nil {
		return err
	}
	path := "."
	generators := []generators.IGenerator{
		schema.NewGenerator(path, d),
		gocode.NewGenerator(path, d),
		proto.NewGenerator(path, d),
		store.NewGenerator(path, d),
		mongo.NewGenerator(path, d),
		gomod.NewGenerator(path, d),
		service.NewGenerator(path, d),
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
