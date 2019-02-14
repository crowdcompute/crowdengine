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

package rpc

import (
	"context"
	"time"

	"github.com/crowdcompute/crowdengine/accounts/keystore"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/log"
	"github.com/crowdcompute/crowdengine/p2p"
)

// AccountsAPI ...
type AccountsAPI struct {
	host *p2p.Host
	ks   *keystore.KeyStore
}

// NewAccountsAPI creates a new accounts API
func NewAccountsAPI(h *p2p.Host, ks *keystore.KeyStore) *AccountsAPI {
	return &AccountsAPI{
		host: h,
		ks:   ks,
	}
}

// CreateAccount creates a new account
func (api *AccountsAPI) CreateAccount(ctx context.Context, passphrase string) (string, string, error) {
	acc, err := api.ks.NewAccount(passphrase)
	common.FatalIfErr(err, "There was an error creating the account")
	log.Printf("The account has been created successfully to the file: {%s}\n", acc.Path)
	return acc.Address, acc.Path, nil
}

// UnlockAccount unlocks an account and returns its token
func (api *AccountsAPI) UnlockAccount(ctx context.Context, accAddr, path, passphrase string) (string, error) {
	// First issue a token
	rawToken, err := api.ks.IssueTokenForAccount(accAddr, keystore.NewTokenClaims("", ""))
	if err != nil {
		return "", err
	}
	// Then unlock the account if there is no issue with the Token creation above
	if err := api.ks.TimedUnlock(accAddr, passphrase, common.TokenTimeout); err != nil {
		return "", err
	}
	toMinutes := float64(common.TokenTimeout) / float64(time.Minute)
	log.Printf("The account {%s} has been unlocked for %.2f minutes... This is the token: {%s} \n", accAddr, toMinutes, rawToken)
	return rawToken, err
}

// ListAccounts creates a new account
func (api *AccountsAPI) ListAccounts(ctx context.Context) []string {
	accounts := api.ks.GetAccounts()
	accAddresses := make([]string, 0)
	for _, acc := range accounts {
		accAddresses = append(accAddresses, acc.Address)
	}
	log.Printf("Here is the list of all accounts of this node: {%s}\n", accAddresses)
	return accAddresses
}

// DeleteAccount deletes the account with the address accAddr and the passphrase passphrase
func (api *AccountsAPI) DeleteAccount(ctx context.Context, accAddr, passphrase string) error {
	err := api.ks.Delete(accAddr, passphrase)
	if err != nil {
		return err
	}
	log.Printf("Deleted this account : {%s}\n", accAddr)
	return nil
}

// LockAccount locks an account
func (api *AccountsAPI) LockAccount(ctx context.Context, accAddr string) error {
	err := api.ks.Lock(accAddr)
	if err != nil {
		log.Println("There was an error locking the account. Error: ", err)
		return err
	}
	log.Printf("The account {%s} has been locked... \n", accAddr)
	return nil
}
