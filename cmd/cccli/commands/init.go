package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/crowdcompute/crowdengine/cmd/utils"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/node"
	"github.com/urfave/cli"
)

var (
	// InitCommand initializes and starts a node
	InitCommand = cli.Command{
		Action:    initNode,
		Name:      "init",
		Usage:     "Initialize a Crowd Compute node",
		ArgsUsage: "<bootnodes>",
		Flags: []cli.Flag{
			utils.P2pPortFlag,
			utils.IPFlag,
			utils.RPCFlag,
		},
		Category: "CC CLI",
		Description: `
		The init command initializes a new node for the network.
		It expects the <bootnodes> file as an argument.`,
	}
)

// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func initNode(ctx *cli.Context) error {
	bootnodesFile := ctx.Args().First()
	port := ctx.Int("port")
	ip := ctx.String("ip")
	rpcFlag := ctx.Bool("rpc") // BoolTFlag returns false if no flag was set

	bootnodes := getBootnodesFromFile(bootnodesFile)

	// Creating a new Node.
	// Takes port and ip command line arguments as parameters
	n, err := node.NewNode(port, ip, bootnodes)
	common.CheckErr(err, "[initNode] Failed to create node.")

	// Create a cancellable context for our GRPC call
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()
	n.Start(ctxt, rpcFlag)

	return nil
}

// getBootnodesFromFile returns a slice of bootnodes given the file bootnodesFile
func getBootnodesFromFile(bootnodesFile string) []string {
	bootnodes := make([]string, 0)
	if bootnodesFile != "" {
		file, err := os.Open(bootnodesFile)
		common.CheckErr(err, "[getBootnodes] Failed to open the bootnodes file.")

		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			bootnodes = append(bootnodes, scanner.Text())
			fmt.Println(scanner.Text())
		}
	}
	return bootnodes
}
