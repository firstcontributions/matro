package commands

import (
	"flag"
	"io/fs"
	"os"

	"github.com/firstcontributions/matro/internal/generators/types"
	"github.com/firstcontributions/matro/internal/parser"
	"github.com/sirupsen/logrus"
)

// CodeGenerator will act as a base class for all code generator commands
type CodeGenerator struct {
	*CommandWriter
	flags    flag.FlagSet
	filepath string
	verbose  bool
	help     bool
	fs.FS
	outputPath string
}

// NewCodeGenerator return a new instance of CodeGenerator
func NewCodeGenerator(writer *CommandWriter) *CodeGenerator {
	return &CodeGenerator{
		CommandWriter: writer,
		FS:            os.DirFS("."),
		outputPath:    ".",
	}
}

// InitFlags will initialize all flags
func (c *CodeGenerator) InitFlags() {
	c.flags.StringVar(&c.filepath, "f", "matro.json", "file path")
	c.flags.StringVar(&c.filepath, "file", "matro.json", "file path")
	c.flags.BoolVar(&c.help, "h", false, "help")
	c.flags.BoolVar(&c.help, "help", false, "help")
	c.flags.BoolVar(&c.verbose, "vv", false, "verbose")
}

// ParseFlags will parse given flags
func (c *CodeGenerator) ParseFlags(args []string) {
	c.flags.Parse(args)
}

// Help prints the help message
func (c *CodeGenerator) Help() {
	helpText := `
	matro <code generator>  -f [--file] <file path>
	[-vv] for verbose

	`
	c.Write(helpText)
}

// GetDefenitionsAndTypes will parse given config and return definitions and typedefs
func (c *CodeGenerator) GetDefenitionsAndTypes() (*parser.Definition, *types.TypeDefs, error) {
	matroConfigs, err := c.FS.Open(c.filepath)
	if err != nil {
		return nil, nil, err
	}
	d, err := parser.NewDefinition().ParseFrom(matroConfigs)
	if err != nil {
		return nil, nil, err
	}
	typeDefs, err := types.GetTypeDefs(d)
	if err != nil {
		return nil, nil, err
	}
	return d, typeDefs, nil
}

// Setup will setup log levels and help command
func (c *CodeGenerator) Setup() bool {
	if c.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
	}
	if c.help {
		c.Help()
		return false
	}
	return true
}
