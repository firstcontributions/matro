package main

import (
	"flag"

	"github.com/firstcontributions/matro/internal/commands"
	"github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	cmd := commands.GetCmd(flag.Args())
	if err := cmd.Exec(); err != nil {
		logrus.Fatal(err)
	}
}
