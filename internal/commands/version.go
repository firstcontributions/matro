package commands

import "fmt"

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
}

// NewVersionCmd returns a new instalce of VersionCmd
func NewVersionCmd() *VersionCmd {
	return &VersionCmd{}
}

//InitFlags initialises the flags if any
func (VersionCmd) InitFlags() {}

// ParseFlags will parse given flags
func (VersionCmd) ParseFlags(args []string) {}

// Help prints VersionCmd text for the command, not needed here
func (VersionCmd) Help() {
	helpText := `
	version
	matro version 
	
	`
	fmt.Println(helpText)
}

// Exec will print the build versionCmd details
func (VersionCmd) Exec() error {
	fmt.Println("Version :", Version)
	fmt.Println("MinVersion :", MinVersion)
	fmt.Println("BuildTime :", BuildTime)
	return nil
}
