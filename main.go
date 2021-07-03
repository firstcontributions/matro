package main

import (
	"github.com/firstcontributions/matro/internal/cleaner"
	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/parser"
)

func main() {
	d, err := parser.NewDefinition().ParseFromFile("./input.json")
	if err != nil {
		panic(err)
	}
	basePath := "./__generated"
	cleaner.Clean(basePath)
	s := generators.GetGenerator(basePath, "schema", d)
	err = s.Generate()
	if err != nil {
		panic(err)
	}
	gc := generators.GetGenerator(basePath, "gocode", d)
	err = gc.Generate()
	if err != nil {
		panic(err)
	}

	gp := generators.GetGenerator(basePath, "proto", d)
	err = gp.Generate()
	if err != nil {
		panic(err)
	}
}
