package commands

import (
	"fmt"

	"github.com/crowdcompute/crowdengine/accounts/keystore"
	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/cmd/terminal"
	"github.com/crowdcompute/crowdengine/common"
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
				Description: `
Print a short summary of all accounts`,
			},
			{
				Name:   "unlock",
				Usage:  "Unlock an existing account",
				Action: UnlockAccount,
				Description: `
Print a short summary of all accounts`,
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
	pass, err := terminal.Stdin.GetPassphrase("Please give a password and not forget this password.", true)
	common.FatalIfErr(err, "Error reading passphrase from terminal")
	cfg := config.GetConfig(ctx)
	ks := keystore.NewKeyStore(cfg.Global.KeystoreDir)
	acc, err := ks.NewAccount(pass)
	fmt.Printf("Your account has been created successfully to the file: {%s}\n", acc.Path)
	return nil
}

// ListAccounts creates a new account for the user
func ListAccounts(ctx *cli.Context) error {
	cfg := config.GetConfig(ctx)
	accounts, err := keystore.GetKeystoreFiles(cfg.Global.KeystoreDir)
	common.FatalIfErr(err, "Unable to get keystore files")
	for i, account := range accounts {
		fmt.Printf("Account #%d: %s\n", i, account)
	}
	return nil
}

// LockAccount locks an existing account
func LockAccount(ctx *cli.Context) error {

	return nil
}

// UnlockAccount unlocks an existing account
func UnlockAccount(ctx *cli.Context) error {
	// fmt.Printf("Here is your token. Use this for further calls: {%s}\n", acc.Token.Raw)
	return nil
}
