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
	"fmt"

	ccsdk "github.com/crowdcompute/cc-go-sdk"
	"github.com/crowdcompute/crowdengine/cmd/ccpush/config"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/urfave/cli"
)

var (
	// ImageCommand is a command for managing images
	ImageCommand = cli.Command{
		Name:     "image",
		Usage:    "Manage running images",
		Category: "Image",
		Description: `
					Manage images. Send images to nodes on the network to run`,
		Subcommands: []cli.Command{
			{
				Name:   "deploy",
				Usage:  "deploy <account> <passphrase> <imgpath> <libp2pID>",
				Action: RunImageOnNode,
				Flags: []cli.Flag{
					config.RPCAddrFlag,
					config.FileserverFlag,
					config.AccAddrFlag,
					config.AccPassphraseFlag,
					config.ImgPathFlag,
					config.Libp2pIDFlag,
				},
				Description: `
				Executes an image to a specified node`,
			},
		},
	}
)

// RunImageOnNode run image on a node
func RunImageOnNode(ctx *cli.Context) error {
	// Check for 3 because help flag is there by default
	if len(ctx.Command.VisibleFlags()) != 7 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
	// Get the clients to communicate with the node
	rpcaddr := ctx.String(config.RPCAddrFlag.Name)
	c := ccsdk.NewCCClient(rpcaddr)
	fileserveraddr := ctx.String(config.FileserverFlag.Name)
	uploadClient := ccsdk.NewUploadClient(fileserveraddr)
	// Get the rest of the flags
	accAddr := ctx.String(config.AccAddrFlag.Name)
	passphrase := ctx.String(config.AccPassphraseFlag.Name)
	imagePath := ctx.String(config.ImgPathFlag.Name)
	libp2pID := ctx.String(config.Libp2pIDFlag.Name)

	// Unlock it
	token, err := c.UnlockAccount(accAddr, passphrase)
	common.FatalIfErr(err, "Couldn't unlock account.")

	// Upload an image to dev node
	imgHash, err := uploadClient.UploadFile(imagePath, token)
	common.FatalIfErr(err, "Couldn't upload image file to dev node.")
	fmt.Println("Result of imgHash is this: ", imgHash)

	// Load image to a node's docker engine
	imgID, err := c.LoadImageToNode(libp2pID, imgHash, token)
	common.FatalIfErr(err, "Couldn't load image node.")
	fmt.Println("Result of LoadImageToNode is this: ", imgID)

	// Execute image
	result, err := c.ExecuteImage(libp2pID, imgID)
	common.FatalIfErr(err, "Couldn't run image to node.")
	fmt.Println("Result of ExecuteImage is this: ", result)
	return nil
}
