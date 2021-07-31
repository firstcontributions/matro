package gomod

import (
	"context"
	"os"
	"os/exec"

	"github.com/firstcontributions/matro/internal/generators/utils"
	"github.com/firstcontributions/matro/internal/parser"
)

// Generator implements graphql server code generator
type Generator struct {
	Repo string
	Path string
}

// NewGenerator returns an instance of graphql server code generator
func NewGenerator(path string, d *parser.Definition) *Generator {
	return &Generator{
		Path: path,
		Repo: d.Repo,
	}
}

// Generate generates a go mod file and get all dependencies
func (g *Generator) Generate(ctx context.Context) error {
	if err := utils.EnsurePath(g.Path); err != nil {
		return err
	}
	if err := g.initGoMod(); err != nil {
		return err
	}
	return g.goGet()
}

func (g *Generator) initGoMod() error {
	cmdGoMod := exec.Command("go", "mod", "init", g.Repo)
	cmdGoMod.Dir = g.Path
	cmdGoMod.Stdout = os.Stdout
	cmdGoMod.Stderr = os.Stderr
	return cmdGoMod.Run()
}

func (g *Generator) goGet() error {
	cmd := exec.Command("go", "get", "-d", "-v", "./...")
	cmd.Dir = g.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
