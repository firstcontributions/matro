package main

import (
	"github.com/firstcontributions/matro/internal/generators"
	"github.com/firstcontributions/matro/internal/parser"
)

func main() {
	d, err := parser.NewDefinition().ParseFromFile("./input.json")
	if err != nil {
		panic(err)
	}
	s := generators.GetGenerator("./__generated", "schema", d)
	err = s.Generate()
	if err != nil {
		panic(err)
	}
	gc := generators.GetGenerator("./__generated", "gocode", d)
	err = gc.Generate()
	if err != nil {
		panic(err)
	}
}
