package commands

import "os"

// Command is the interface implemented by types
type Command interface {
	InitFlags()
	ParseFlags([]string)
	Help()
	Exec() error
}

var (
	server     *Server
	relay      *Relay
	helpCmd    *Help
	versionCmd *VersionCmd
)

// initializes when package initializes
func init() {
	w := NewCommandWriter(os.Stdout)

	server = NewServer(w)
	server.InitFlags()

	relay = NewRelay(w)
	relay.InitFlags()

	helpCmd = NewHelp(w)
	helpCmd.InitFlags()

	versionCmd = NewVersionCmd(w)
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
