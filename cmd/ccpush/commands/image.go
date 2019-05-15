package commands

import (
	"fmt"

	ccsdk "github.com/crowdcompute/cc-go-sdk"
	"github.com/crowdcompute/crowdengine/cmd/ccpush/config"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/urfave/cli"
)

var (
	c            = ccsdk.NewCCClient("http://localhost:8085")
	uploadClient = ccsdk.NewUploadClient("http://localhost:8085/upload")
)

var (
	// ImageCommand is a command for managing images
	ImageCommand = cli.Command{
		Name:     "image",
		Usage:    "Manage running images",
		Category: "Images",
		Description: `
					Manage images. Send images to nodes on the network to run`,
		Subcommands: []cli.Command{
			{
				Name:   "deploy",
				Usage:  "deploy <account> <passphrase> <imgpath> <libp2pID>",
				Action: RunImageOnNode,
				Flags: []cli.Flag{
					config.AccAddrFlag,
					config.AccPassphraseFlag,
					config.ImgPathFlag,
					config.Libp2pIDFlag,
				},
				Description: `
				Executes an image to a specified node`,
			},
			// {
			// 	Name:   "list",
			// 	Usage:  "list <libp2pID> <token>",
			// 	Action: ListNodeContainers,
			// 	Flags: []cli.Flag{
			// 		config.Libp2pIDFlag,
			// 		config.TokenFlag,
			// 	},
			// 	Description: `
			// 	Lists containers of a specified node`,
			// },
		},
	}
)

// RunImageOnNode run image on a node
func RunImageOnNode(ctx *cli.Context) error {
	// Check for 3 because help flag is there by default
	if len(ctx.Command.VisibleFlags()) != 5 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
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

// ListNodeContainers lists containers of a specified node
// func ListNodeContainers(ctx *cli.Context) error {
// 	// Check for 3 because help flag is there by default
// 	if len(ctx.Command.VisibleFlags()) != 3 {
// 		return fmt.Errorf("Please give account and passphrase flags")
// 	}
// 	libp2pID := ctx.String(config.Libp2pIDFlag.Name)
// 	token := ctx.String(config.TokenFlag.Name)

// 	// List containers
// 	contList, err := c.ListNodeContainers(libp2pID, token)
// 	if err != nil {
// 		fmt.Println("Couldn't list image containers.", err)
// 	}
// 	fmt.Println("Result of ListNodeContainers is this: ", contList)
// 	return nil
// }
