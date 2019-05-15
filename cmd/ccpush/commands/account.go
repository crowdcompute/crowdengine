package commands

import (
	"fmt"

	"github.com/crowdcompute/crowdengine/cmd/ccpush/config"
	"github.com/urfave/cli"
)

var (
	// AccountCommand is a commnad for managing accounts
	AccountCommand = cli.Command{
		Name:     "account",
		Usage:    "Manage accounts",
		Category: "Accounts",
		Description: `
					Manage accounts, create, lock, unlock and import existing ones`,
		Subcommands: []cli.Command{
			{
				Name:   "new",
				Usage:  "add a new account",
				Action: NewAccount,
			},
			{
				Name:   "list",
				Usage:  "Print summary of existing accounts",
				Action: ListAccounts,
				Description: `
							Print a short summary of all accounts`,
			},
			{
				Name:   "lock",
				Usage:  "Locks an existing account",
				Action: LockAccount,
				Flags: []cli.Flag{
					config.AccAddrFlag,
				},
				Description: `
							Locks a specific account`,
			},
			{
				Name:   "unlock",
				Usage:  "Unlock an existing account",
				Action: UnlockAccount,
				Flags: []cli.Flag{
					config.AccAddrFlag,
					config.AccPassphraseFlag,
				},
				Description: `
							Unlocks a specific account`,
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
func NewAccount(ctx *cli.Context) error {
	fmt.Println("Newaccount command to be implemented")
	return nil
}

// ListAccounts lists all accounts of the node
func ListAccounts(ctx *cli.Context) error {
	accounts, err := c.ListAccounts()
	if err != nil {
		fmt.Println("Unable to get accounts", err)
	}
	for i, account := range accounts {
		fmt.Printf("Account #%d: %s\n", i, account)
	}
	return nil
}

// LockAccount locks an existing account
func LockAccount(ctx *cli.Context) error {
	fmt.Println("LockAccount command to be implemented")
	return nil
}

// UnlockAccount unlocks an existing account
func UnlockAccount(ctx *cli.Context) error {
	// Check for 3 because help flag is there by default
	if len(ctx.Command.VisibleFlags()) != 3 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
	accAddr := ctx.String(config.AccAddrFlag.Name)
	passphrase := ctx.String(config.AccPassphraseFlag.Name)
	// Unlock it
	token, err := c.UnlockAccount(accAddr, passphrase)
	if err != nil {
		fmt.Println("Can't issue token or unlock account.", err)
	}
	fmt.Println("token: ", token)
	return nil
}
