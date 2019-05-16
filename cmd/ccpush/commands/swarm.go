package commands

import (
	"encoding/json"
	"fmt"

	ccsdk "github.com/crowdcompute/cc-go-sdk"
	"github.com/crowdcompute/crowdengine/cmd/ccpush/config"
	"github.com/urfave/cli"
)

var (
	// SwarmCommand is a command for managing swarms
	SwarmCommand = cli.Command{
		Name:     "swarm",
		Usage:    "Manage running swarm",
		Category: "Swarm",
		Description: `
					Manage swarm on nodes.`,
		Subcommands: []cli.Command{
			{
				Name:   "deploy",
				Usage:  "deploy <account> <passphrase> <imgpath> <libp2pID>",
				Action: createAndRunSwarm,
				Flags: []cli.Flag{
					config.RPCAddrFlag,
					// config.AccAddrFlag,
					// config.AccPassphraseFlag,
					// config.ImgPathFlag,
					config.Libp2pIDFlag,
				},
				Description: `
				Executes images as part of a docker swarm`,
			},
			{
				Name:   "stop",
				Usage:  "stop <account> <passphrase> <imgpath> <libp2pID>",
				Action: createAndRunSwarm,
				Flags: []cli.Flag{
					config.RPCAddrFlag,
					config.FileserverFlag,
					// config.AccAddrFlag,
					// config.AccPassphraseFlag,
					// config.ImgPathFlag,
					config.Libp2pIDFlag,
				},
				Description: `
				Executes images as part of a docker swarm`,
			},
		},
	}
)

// RunImageOnNode run image on a node
func createAndRunSwarm(ctx *cli.Context) error {
	// Check for 3 because help flag is there by default
	if len(ctx.Command.VisibleFlags()) != 3 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
	// Get the client to communicate with the node
	rpcaddr := ctx.String(config.RPCAddrFlag.Name)
	c := ccsdk.NewCCClient(rpcaddr)

	// accAddr := ctx.String(config.AccAddrFlag.Name)
	// passphrase := ctx.String(config.AccPassphraseFlag.Name)
	// imagePath := ctx.String(config.ImgPathFlag.Name)
	libp2pID := ctx.String(config.Libp2pIDFlag.Name)

	// Unlock it
	// token, err := c.UnlockAccount(accAddr, passphrase)
	// common.FatalIfErr(err, "Couldn't unlock account.")
	// Create and run swarm service
	type swarmTask struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	task := swarmTask{"JustATag", "animage"}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = c.RunSwarmService(string(taskBytes), []string{libp2pID})
	if err != nil {
		fmt.Println("There was an error running the docker swarm: ", err)
	}
	return nil
}
