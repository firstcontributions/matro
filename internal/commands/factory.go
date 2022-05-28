package commands

// Command is the interface implemented by types
type Command interface {
	InitFlags()
	ParseFlags([]string)
	Help()
	Exec() error
}

var (
	server     = NewServer()
	relay      = NewRelay()
	helpCmd    = NewHelp()
	versionCmd = NewVersionCmd()
)

// initializes when package initializes
func init() {
	server.InitFlags()
	relay.InitFlags()
	helpCmd.InitFlags()
	versionCmd.InitFlags()
}

// GetCmd will get command by flags
func GetCmd(args []string) Command {
	if len(args) == 0 {
		return helpCmd
	}
	switch args[0] {
	case "server":
		server.ParseFlags(args[1:])
		return server
	case "relay":
		relay.ParseFlags(args[1:])
		return relay
	case "version":
		return versionCmd
	default:
		return helpCmd
	}
}
