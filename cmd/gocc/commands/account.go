package commands

import (
	"fmt"

	"github.com/crowdcompute/crowdengine/accounts/keystore"
	"github.com/urfave/cli"
)

var (
	// AccountCommand adds accounting functionlaity
	AccountCommand = cli.Command{
		Name:     "account",
		Usage:    "Manage accounts",
		Category: "Accounts",
		Description: `
Manage accounts, create update and import new stuff`,
		Subcommands: []cli.Command{
			{
				Name:   "new",
				Usage:  "add a new account",
				Action: NewAccount,
			},
			{
				Name:  "import",
				Usage: "import an existing account",
				Action: func(c *cli.Context) error {
					fmt.Println("importing account...")
					// geenral error
					//return cli.NewExitError("general error", 1)
					return nil
				},
			},
		},
	}
)

// NewAccount creates a new account for the user
func NewAccount(c *cli.Context) error {
	_, fileName := keystore.NewKeyAndStoreToFile()
	fmt.Printf("Your account has been created successfully to the file: {%s}\n", fileName)
	return nil
}
