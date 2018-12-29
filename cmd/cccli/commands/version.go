package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

var (
	// VersionCommand prints the version of the program
	VersionCommand = cli.Command{
		Action:    version,
		Name:      "version",
		Usage:     "Version",
		ArgsUsage: "<>",
		Category:  "CC CLI",
		Description: `
		Version`,
	}
)

// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func version(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	// print Version
	return nil
}
