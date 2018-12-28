package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

var (
	TaskRootCommand = cli.Command{
		Name:      "task",
		Usage:     "Task management",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// 	utils.DataDirFlag,
		// 	utils.LightModeFlag,
		// },
		Category: "CC CLI",
		Description: `
		Task`,
	}
	taskListSubCommand = cli.Command{
		Action:    taskList,
		Name:      "list",
		Usage:     "list [addr]",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// 	utils.DataDirFlag,
		// 	utils.LightModeFlag,
		// },
		Description: `
		Show active tasks`,
	}
	taskStartSubCommand = cli.Command{
		Action:    taskStart,
		Name:      "start",
		Usage:     "start <task.yaml>",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// 	utils.DataDirFlag,
		// 	utils.LightModeFlag,
		// },
		Description: `
		Start task`,
	}
	taskStatusSubCommand = cli.Command{
		Action:    taskStatus,
		Name:      "status",
		Usage:     "status <addr> <task_id>",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// 	utils.DataDirFlag,
		// 	utils.LightModeFlag,
		// },
		Description: `
		Show the task's status`,
	}
	taskLogsSubCommand = cli.Command{
		Action:    taskLogs,
		Name:      "logs",
		Usage:     "logs <addr> <task_id>",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// 	utils.DataDirFlag,
		// 	utils.LightModeFlag,
		// },
		Description: `
		Show the task's logs`,
	}
	taskStopSubCommand = cli.Command{
		Action:    taskStop,
		Name:      "stop",
		Usage:     "stop <addr> <task_id>",
		ArgsUsage: "<>",
		// Flags: []cli.Flag{
		// 	utils.DataDirFlag,
		// 	utils.LightModeFlag,
		// },
		Description: `
		Stop a task`,
	}
)

func init() {
	TaskRootCommand.Subcommands = []cli.Command{
		taskListSubCommand,
		taskStartSubCommand,
		taskStatusSubCommand,
		taskLogsSubCommand,
		taskStopSubCommand,
	}
}

func taskList(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	// node := makeFullNode(ctx)
	// startNode(ctx, node)
	// node.Wait()
	return nil
}

func taskStart(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	// node := makeFullNode(ctx)
	// startNode(ctx, node)
	// node.Wait()
	return nil
}

func taskStatus(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	// node := makeFullNode(ctx)
	// startNode(ctx, node)
	// node.Wait()
	return nil
}

func taskLogs(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	// node := makeFullNode(ctx)
	// startNode(ctx, node)
	// node.Wait()
	return nil
}

func taskStop(ctx *cli.Context) error {
	fmt.Println("To be implemented...")
	// node := makeFullNode(ctx)
	// startNode(ctx, node)
	// node.Wait()
	return nil
}
