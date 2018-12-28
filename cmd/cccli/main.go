package main

import (
	"os"
	"sort"

	cmd "github.com/crowdcompute/crowdengine/cmd/cccli/commands"
	utils "github.com/crowdcompute/crowdengine/cmd/utils"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/urfave/cli"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	app       = utils.NewApp(gitCommit, "the Crowd-compute command line interface")
)

// Initialize the CLI app
func init() {
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2013-2017 The Crowd Compute Authors"
	app.Commands = []cli.Command{
		cmd.InitCommand,
		cmd.WorkerRootCommand,
		cmd.TaskRootCommand,
		cmd.ServiceRootCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.After = func(ctx *cli.Context) error {
		// debug.Exit()
		// console.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
	err := app.Run(os.Args)
	common.CheckErr(err, "[main] Failed to run the app.")
}
