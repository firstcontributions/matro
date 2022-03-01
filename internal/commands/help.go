package commands

import "fmt"

// Help is the help command
type Help struct {
}

// NewHelp returns a new instalce of Help
func NewHelp() *Help {
	return &Help{}
}

//InitFlags initialises the flags if any
func (Help) InitFlags() {}

// ParseFlags will parse given flags
func (Help) ParseFlags(args []string) {}

// Help prints help text for the command, not needed here
func (Help) Help() {
	txt := `
	matro help
	Help command prints help
	`
	fmt.Println(txt)
}

// Exec will print help text for envparser
func (h Help) Exec() error {
	fmt.Println("Usage")
	h.Help()
	NewServer().Help()
	NewRelay().Help()
	NewVersionCmd().Help()
	return nil
}
