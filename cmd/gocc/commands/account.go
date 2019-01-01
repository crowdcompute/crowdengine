package commands

import (
	"fmt"

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
				Name:  "new",
				Usage: "add a new account",
				Action: func(c *cli.Context) error {
					fmt.Println("new account...")
					return nil
				},
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
