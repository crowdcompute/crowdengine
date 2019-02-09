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
func (api *AccountsAPI) CreateAccount(ctx context.Context, passphrase string) (string, error) {
	acc, err := api.ks.NewAccount(passphrase)
	common.FatalIfErr(err, "There was an error creating the account")
	log.Printf("The account has been created successfully to the file: {%s}\n", acc.Path)
	return acc.Path, nil
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
	return rawToken, err
}

// LockAccount locks an account
func (api *AccountsAPI) LockAccount(ctx context.Context, accAddr string) (peers []string) {
	api.ks.Lock(accAddr)
	return
}
