// Copyright 2018 The crowdcompute:crowdengine Authors
// This file is part of the crowdcompute:crowdengine library.
//
// The crowdcompute:crowdengine library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The crowdcompute:crowdengine library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the crowdcompute:crowdengine library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/cmd"
	"github.com/crowdcompute/crowdengine/cmd/gocc/commands"
	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/node"
	"github.com/urfave/cli"
)

var (
	// Version is passed using the make file
	Version string

	// GitCommit is used to reference the commit used for the build
	GitCommit string

	// App is an instance of a cli app
	App = cmd.NewApp(GitCommit)
)

func init() {
	// App.HideVersion = true
	App.Action = gocc
	App.Version = Version
	App.Commands = []cli.Command{
		commands.AccountCommand,
	}
	App.Flags = config.GOCCAppFlags
	sort.Sort(cli.CommandsByName(App.Commands))
	App.After = func(ctx *cli.Context) error {
		// debug.Exit()
		// console.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func gocc(ctx *cli.Context) error {
	// create default config
	cfg := config.DefaultConfig()

	// if config file is given, load it
	confFile := ctx.String("config")
	if confFile != "" {
		config.LoadTomlConfig(ctx, cfg)
	}

	// apply flags to config
	config.ApplyFlags(ctx, cfg)

	// create and start node
	if node, err := node.NewNode(cfg); err != nil {
		log.Fatal(err)
	} else {
		node.Start(ctx)
	}

	return nil
}

func main() {
	if err := App.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
