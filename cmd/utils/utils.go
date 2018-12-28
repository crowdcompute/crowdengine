package utils

import (
	"github.com/urfave/cli"
)

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, usage string) *cli.App {
	app := cli.NewApp()
	app.Name = "test-cli"
	app.Author = ""
	//app.Authors = nil
	app.Email = ""
	if gitCommit != "" {
		app.Version += "-" + gitCommit[:8]
	}
	app.Usage = usage
	return app
}

var (
	// Init settings
	P2pPortFlag = cli.IntFlag{
		Name:  "port, p",
		Usage: "The p2p port the node will listen to.",
		Value: 12000,
	}
	IPFlag = cli.StringFlag{
		Name:  "ip",
		Usage: "Setting the IP address",
		Value: "127.0.0.0",
	}
	// Node settings
	RPCFlag = cli.BoolFlag{
		Name:  "rpc",
		Usage: "Setting this flag will make this node listet to JSON-RPC calls",
	}
)
