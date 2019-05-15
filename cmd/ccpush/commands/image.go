package commands

import (
	"fmt"

	ccsdk "github.com/crowdcompute/cc-go-sdk"
	"github.com/crowdcompute/crowdengine/cmd/ccpush/config"
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
				Name:   "upload",
				Usage:  "upload <filePath> <token>",
				Action: UploadImageToDevNode,
				Flags: []cli.Flag{
					config.FilePathFlag,
					config.TokenFlag,
				},
				Description: `uploads an image to the dev node`,
			},
			{
				Name:   "load",
				Usage:  "load <libp2pID> <imgHash> <token>",
				Action: LoadImageToNode,
				Flags: []cli.Flag{
					config.Libp2pIDFlag,
					config.ImgHashFlag,
					config.TokenFlag,
				},
				Description: `
				loads an image to a node's docker engine`,
			},
			{
				Name:   "run",
				Usage:  "run <libp2pID> <imgID>",
				Action: ExecuteImage,
				Flags: []cli.Flag{
					config.Libp2pIDFlag,
					config.ImgIDFlag,
				},
				Description: `
				Executes an image to a specified node`,
			},
			{
				Name:   "list",
				Usage:  "list <libp2pID> <token>",
				Action: ListNodeContainers,
				Flags: []cli.Flag{
					config.Libp2pIDFlag,
					config.TokenFlag,
				},
				Description: `
				Lists containers of a specified node`,
			},
		},
	}
)

// UploadImageToDevNode uploads an image to the dev node
func UploadImageToDevNode(ctx *cli.Context) error {
	// Check for 3 because help flag is there by default
	if len(ctx.Command.VisibleFlags()) != 3 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
	filePath := ctx.String(config.FilePathFlag.Name)
	token := ctx.String(config.TokenFlag.Name)

	// Upload an image to dev node
	imgHash, err := uploadClient.UploadFile(filePath, token)
	if err != nil {
		fmt.Println("Couldn't upload image file to dev node.", err)
	}
	fmt.Println("Result of imgHash is this: ", imgHash)
	return nil
}

// LoadImageToNode loads an image to a node's docker engine
func LoadImageToNode(ctx *cli.Context) error {
	// Check for 4 because help flag is there by default
	if len(ctx.Command.VisibleFlags()) != 4 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
	libp2pID := ctx.String(config.Libp2pIDFlag.Name)
	imgHash := ctx.String(config.ImgHashFlag.Name)
	token := ctx.String(config.TokenFlag.Name)

	// Load image to a node's docker engine
	imgID, err := c.LoadImageToNode(libp2pID, imgHash, token)
	if err != nil {
		fmt.Println("Couldn't load image file to node.", err)
	}
	fmt.Println("Result of LoadImageToNode is this: ", imgID)

	return nil
}

// ExecuteImage executes an image to a specified node
func ExecuteImage(ctx *cli.Context) error {
	// Check for 3 because help flag is there by default
	if len(ctx.Command.VisibleFlags()) != 3 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
	libp2pID := ctx.String(config.Libp2pIDFlag.Name)
	imgID := ctx.String(config.ImgIDFlag.Name)

	// Execute image
	result, err := c.ExecuteImage(libp2pID, imgID)
	if err != nil {
		fmt.Println("Couldn't execute the image to node.", err)
	}
	fmt.Println("Result of ExecuteImage is this: ", result)
	return nil
}

// ListNodeContainers lists containers of a specified node
func ListNodeContainers(ctx *cli.Context) error {
	// Check for 3 because help flag is there by default
	if len(ctx.Command.VisibleFlags()) != 3 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
	libp2pID := ctx.String(config.Libp2pIDFlag.Name)
	token := ctx.String(config.TokenFlag.Name)

	// List containers
	contList, err := c.ListNodeContainers(libp2pID, token)
	if err != nil {
		fmt.Println("Couldn't list image containers.", err)
	}
	fmt.Println("Result of ListNodeContainers is this: ", contList)
	return nil
}
