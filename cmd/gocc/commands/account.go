package commands

import (
	"fmt"
	"log"

	"github.com/crowdcompute/crowdengine/accounts/keystore"
	"github.com/crowdcompute/crowdengine/cmd/terminal"
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
	pass, err := terminal.Stdin.GetPassphrase("Please give a password and not forget this password.", true)
	if err != nil {
		log.Fatalf("Error reading passphrase from terminal: %v", err)
	}
	_, fileName := keystore.NewKeyAndStoreToFile(pass)
	fmt.Printf("Your account has been created successfully to the file: {%s}\n", fileName)
	return nil
}
