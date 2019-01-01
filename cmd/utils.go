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

package cmd

import (
	"github.com/urfave/cli"
)

var (
	// AppHelpTemplate is the main tpl
	AppHelpTemplate = `NAME:
{{.Name}} - {{.Usage}}

Copyright 2017-2018 The crowdcompute dev team

USAGE:
{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
{{if len .Authors}}
AUTHOR:
{{range .Authors}}{{ . }}{{end}}
{{end}}{{if .Version}}
VERSION:
  {{.Version}}
  {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
{{range .VisibleFlags}}{{.}}
{{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
{{end}}
`
)

// NewApp creates an app
func NewApp(gitCommit string) *cli.App {
	app := cli.NewApp()
	//cli.AppHelpTemplate = AppHelpTemplate
	app.CustomAppHelpTemplate = AppHelpTemplate
	app.Name = "gocc"
	app.Usage = "command line interface"
	app.Copyright = "Copyright 2018 The CrowdCompute Authors"
	if gitCommit != "" {
		app.Version += "-" + gitCommit[:8]
	}
	return app
}
