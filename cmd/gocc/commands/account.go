package commands

import (
	"fmt"
	"time"

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
				Description: `
							Locks a specific account`,
			},
			{
				Name:   "unlock",
				Usage:  "Unlock an existing account",
				Action: UnlockAccount,
				Description: `
							Unlocks a specific account`,
				Flags: []cli.Flag{
					config.AccAddrFlag,
					config.AccPassphraseFlag,
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

// ListAccounts lists all accounts of the node
func ListAccounts(ctx *cli.Context) error {
	cfg := config.GetConfig(ctx)
	// Need to pass the keystore path to
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
	// help flag is there as well
	if len(ctx.Command.VisibleFlags()) != 3 {
		return fmt.Errorf("Please give account and passphrase flags")
	}
	accAddr := ctx.String(config.AccAddrFlag.Name)
	passphrase := ctx.String(config.AccPassphraseFlag.Name)

	cfg := config.GetConfig(ctx)
	ks := keystore.NewKeyStore(cfg.Global.KeystoreDir)
	// First issue a token
	rawToken, err := ks.IssueTokenForAccount(accAddr, keystore.NewTokenClaims("", ""))
	if err != nil {
		fmt.Printf("cant issue token {%s} ", err)
		return err
	}
	// Then unlock the account if there is no issue with the Token creation above
	if err := ks.TimedUnlock(accAddr, passphrase, common.TokenTimeout); err != nil {
		fmt.Printf("cant unlock account {%s} ", err)
		return err
	}
	toMinutes := float64(common.TokenTimeout) / float64(time.Minute)
	fmt.Printf("The account {%s} has been unlocked for %.2f minutes... This is your token: {%s} \n", accAddr, toMinutes, rawToken)
	return nil
}
