package main

import (
	"github.com/firstcontributions/matro/internal/generators/graphql"
	"github.com/firstcontributions/matro/internal/parser"
)

func main() {
	d, err := parser.NewDefinition().ParseFromFile("./input.json")
	if err != nil {
		panic(err)
	}
	s := graphql.GetGenerator("./__generated", "schema", d)
	err = s.Generate()
	if err != nil {
		panic(err)
	}
	gc := graphql.GetGenerator("./__generated", "gocode", d)
	err = gc.Generate()
	if err != nil {
		panic(err)
	}
}
