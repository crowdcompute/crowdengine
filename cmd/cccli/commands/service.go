package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

var (
	// ServiceRootCommand runs a new service to a node
	ServiceRootCommand = cli.Command{
		Name:      "service",
		Usage:     "Service management",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// },
		Category: "CC CLI",
		Description: `
		Service`,
	}
	serviceCreateSubCommand = cli.Command{
		Action:    serviceCreate,
		Name:      "create",
		Usage:     "create [addr]",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// },
		Description: `
		Create a service`,
	}
)

func init() {
	ServiceRootCommand.Subcommands = []cli.Command{
		serviceCreateSubCommand,
	}
}

func serviceCreate(ctx *cli.Context) error {
	fmt.Println("To be implemented")
	return nil
}
