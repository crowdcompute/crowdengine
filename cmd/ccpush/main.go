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

	"github.com/crowdcompute/crowdengine/cmd"
	"github.com/crowdcompute/crowdengine/cmd/ccpush/commands"
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
	App.Action = ccpush
	App.Version = Version
	// App.Flags = config.CCPushFlags
	App.Commands = []cli.Command{
		commands.ImageCommand,
	}
	sort.Sort(cli.CommandsByName(App.Commands))
	App.After = func(ctx *cli.Context) error {
		return nil
	}
}

func ccpush(ctx *cli.Context) error {

	return nil
}

func main() {
	if err := App.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
