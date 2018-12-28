package commands

import (
	"fmt"

	"github.com/crowdcompute/crowdengine/manager"
	"github.com/urfave/cli"
)

var (
	ServiceRootCommand = cli.Command{
		Name:      "service",
		Usage:     "Service management",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// 	utils.DataDirFlag,
		// 	utils.LightModeFlag,
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
		// 	utils.DataDirFlag,
		// 	utils.LightModeFlag,
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
	fmt.Println("Testing...")
	// p2pport := 12000
	// host := makeRandomHost(*p2pport)
	resp, err := manager.GetInstance().InitSwarm("", "127.0.0.1:12000")
	if err != nil {
		fmt.Printf("[serviceCreate] Failed to init swarm: %s", err)
		return err
	}
	fmt.Printf("resp: %s", resp)

	return nil
}
