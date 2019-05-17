package commands

import (
	"encoding/json"
	"fmt"

	ccsdk "github.com/crowdcompute/cc-go-sdk"
	"github.com/crowdcompute/crowdengine/cmd/ccpush/config"
	"github.com/crowdcompute/crowdengine/common"
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
				Usage:  "create a docker service and run it on the specified nodes",
				Action: createAndRunSwarm,
				Flags: []cli.Flag{
					config.RPCAddrFlag,
					config.Libp2pIDFlag,
				},
				Description: `
				Executes images as part of a docker swarm`,
			},
			{
				Name:   "stop",
				Usage:  "stop <account> <passphrase> <imgpath> <libp2pID>",
				Action: LeaveSwarm,
				Flags: []cli.Flag{
					config.RPCAddrFlag,
					config.Libp2pIDFlag,
				},
				Description: `
				Forces all swarm nodes to leave docker swarm`,
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

	libp2pID := ctx.String(config.Libp2pIDFlag.Name)
	ids := common.CommaStringToSlice(libp2pID)

	type swarmService struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	service := swarmService{"JustATag", "animage"}
	taskBytes, err := json.Marshal(service)
	if err != nil {
		fmt.Println("Error marshaling service: ", err)
	}
	err = c.RunSwarmService(string(taskBytes), ids)
	return err
}

// LeaveSwarm makes all connected swarm nodes leaving the swarm
func LeaveSwarm(ctx *cli.Context) error {
	// Check for 3 because help flag is there by default
	if len(ctx.Command.VisibleFlags()) != 3 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
	// Get the client to communicate with the node
	rpcaddr := ctx.String(config.RPCAddrFlag.Name)
	c := ccsdk.NewCCClient(rpcaddr)

	libp2pID := ctx.String(config.Libp2pIDFlag.Name)
	ids := common.CommaStringToSlice(libp2pID)

	err := c.StopSwarmService(ids)
	return err
}
