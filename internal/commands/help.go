package commands

// Help is the help command
type Help struct {
	*CommandWriter
}

// NewHelp returns a new instalce of Help
func NewHelp(writer *CommandWriter) *Help {
	return &Help{
		CommandWriter: writer,
	}
}

//InitFlags initialises the flags if any
func (Help) InitFlags() {}

// ParseFlags will parse given flags
func (Help) ParseFlags(args []string) {}

// Help prints help text for the command, not needed here
func (h Help) Help() {
	txt := `
	matro help
	Help command prints help

	`
	h.Write(txt)
}

// Exec will print help text for envparser
func (h Help) Exec() error {
	h.Write("Usage\n")
	h.Help()
	NewServer(h.CommandWriter).Help()
	NewRelay(h.CommandWriter).Help()
	NewVersionCmd(h.CommandWriter).Help()
	return nil
}
