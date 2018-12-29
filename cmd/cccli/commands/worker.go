package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

var (
	// WorkerRootCommand runs a new task to a node
	WorkerRootCommand = cli.Command{
		Action:    workerManagement,
		Name:      "worker",
		Usage:     "Worker management",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// },
		Category: "CC CLI",
		Description: `
		Worker`,
	}
	workerStatusSubCommand = cli.Command{
		Action:    workerStatus,
		Name:      "status",
		Usage:     "Worker management",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// },
		Description: `
		Worker status`,
	}
	workerTasksSubCommand = cli.Command{
		Action:    workerTasks,
		Name:      "tasks",
		Usage:     "Worker tasks",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// },
		Description: `
		Worker tasks`,
	}
	workerDevicesSubCommand = cli.Command{
		Action:    workerDevices,
		Name:      "devices",
		Usage:     "Worker devices",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// },
		Description: `
		Worker devices`,
	}
)

func init() {
	WorkerRootCommand.Subcommands = []cli.Command{
		workerStatusSubCommand,
		workerTasksSubCommand,
		workerDevicesSubCommand,
	}
}

// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func workerManagement(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	return nil
}

func workerStatus(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	return nil
}

func workerTasks(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	return nil
}

func workerDevices(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	return nil
}
