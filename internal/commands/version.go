package commands

var (
	// Version is the build version eg: v1.0.0
	Version string
	// MinVersion is the latest commit hash on building
	MinVersion string
	// BuildTime is the unix timestamp when the build was done
	BuildTime string
)

// VersionCmd is the VersionCmd command
type VersionCmd struct {
	*CommandWriter
}

// NewVersionCmd returns a new instalce of VersionCmd
func NewVersionCmd(writer *CommandWriter) *VersionCmd {
	return &VersionCmd{
		CommandWriter: writer,
	}
}

//InitFlags initialises the flags if any
func (VersionCmd) InitFlags() {}

// ParseFlags will parse given flags
func (VersionCmd) ParseFlags(args []string) {}

// Help prints VersionCmd text for the command, not needed here
func (v VersionCmd) Help() {
	helpText := `
	version
	matro version 
	
	`
	v.Write(helpText)
}

// Exec will print the build versionCmd details
func (v VersionCmd) Exec() error {
	v.Write("Version :", Version, "\n")
	v.Write("MinVersion :", MinVersion, "\n")
	v.Write("BuildTime :", BuildTime, "\n")
	return nil
}
