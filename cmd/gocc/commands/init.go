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

package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/node"
	"github.com/urfave/cli"
)

var (
	// P2pPortFlag p2p port
	P2pPortFlag = cli.IntFlag{
		Name:  "port, p",
		Usage: "The p2p port the node will listen to.",
		Value: 12000,
	}

	// IPFlag listen
	IPFlag = cli.StringFlag{
		Name:  "ip",
		Usage: "Setting the IP address",
		Value: "127.0.0.0",
	}
	// RPCFlag settings
	RPCFlag = cli.BoolFlag{
		Name:  "rpc",
		Usage: "Setting this flag will make this node listet to JSON-RPC calls",
	}
)

var (
	// InitCommand initializes and starts a node
	InitCommand = cli.Command{
		Action:    initNode,
		Name:      "init",
		Usage:     "Initialize a Crowd Compute node",
		ArgsUsage: "<bootnodes>",
		Flags: []cli.Flag{
			P2pPortFlag,
			IPFlag,
			RPCFlag,
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
