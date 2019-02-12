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

package keystore

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/crowdcompute/crowdengine/common"
)

var (
	passphrase = "test"
)

func tmpKeyStore(t *testing.T) (string, *KeyStore) {
	dir, err := ioutil.TempDir("", "ccompute-keystore-test")
	if err != nil {
		t.Fatal(err)
	}
	return dir, NewKeyStore(dir)
}

func TestKeyStore(t *testing.T) {
	dir, ks := tmpKeyStore(t)
	defer os.RemoveAll(dir)

	a, err := ks.NewAccount(passphrase)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(a.Path, dir) {
		t.Errorf("account file %s doesn't have dir prefix", a.Path)
	}
	if _, err := ks.Find(a.Address); err != nil {
		t.Errorf("HasAccount(%x). Account should have been there", a.Address)
	}
	if err := ks.Delete(a.Address, passphrase); err != nil {
		t.Errorf("Delete error: %v", err)
	}
	if common.FileExist(a.Path) {
		t.Errorf("account file %s should be gone after Delete", a.Path)
	}
	if _, err := ks.Find(a.Address); err != ErrCouldNotFindAccount {
		t.Errorf("Find(%x). Account shouldn't been there", a.Address)
	}
}

func TestGetAccounts(t *testing.T) {
	dir, ks := tmpKeyStore(t)
	defer os.RemoveAll(dir)

	if accounts := ks.GetAccounts(); len(accounts) != 0 {
		t.Errorf("Accounts length should be zero. Got %d", len(accounts))
	}
	a, err := ks.NewAccount(passphrase)
	if err != nil {
		t.Errorf("There was an error creating new account: %s", err)
	}
	if accounts := ks.GetAccounts(); len(accounts) != 1 {
		t.Errorf("Accounts length should be one. Got %d", len(accounts))
	}
	if err := ks.Delete(a.Address, passphrase); err != nil {
		t.Errorf("Delete error: %v", err)
	}
	if accounts := ks.GetAccounts(); len(accounts) != 0 {
		t.Errorf("Accounts length should be zero. Got %d", len(accounts))
	}
}

func TestTimedUnlock(t *testing.T) {
	dir, ks := tmpKeyStore(t)
	defer os.RemoveAll(dir)

	a, err := ks.NewAccount(passphrase)
	if err != nil {
		t.Fatal(err)
	}
	accAddr := a.Address
	var rawToken string
	// We have to issue a token first
	if rawToken, err = ks.IssueTokenForAccount(accAddr, NewTokenClaims("", "")); err != nil {
		t.Fatal(err)
	}
	// Unlocking the account
	if err = ks.TimedUnlock(accAddr, passphrase, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if _, err := ks.GetKeyIfUnlockedAndValid(rawToken); err != nil {
		t.Fatal("Account shouldn't have been locked or have any problems returning a valid key")
	}

	// Expiring lock
	time.Sleep(250 * time.Millisecond)
	if _, err := ks.GetKeyIfUnlockedAndValid(rawToken); err != ErrLocked {
		t.Fatal("Account should have been locked")
	}
}

func TestOverrideUnlock(t *testing.T) {
	dir, ks := tmpKeyStore(t)
	defer os.RemoveAll(dir)

	a, err := ks.NewAccount(passphrase)
	if err != nil {
		t.Fatal(err)
	}
	accAddr := a.Address
	var rawToken string
	// We have to issue a token first
	if rawToken, err = ks.IssueTokenForAccount(accAddr, NewTokenClaims("", "")); err != nil {
		t.Fatal(err)
	}

	// Unlock for a period of time
	if err = ks.TimedUnlock(accAddr, passphrase, 5*time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, err := ks.GetKeyIfUnlockedAndValid(rawToken); err != nil {
		t.Fatal("Account shouldn't have been locked or have any problems returning a valid key")
	}

	// reset unlock to a shorter period, invalidates the previous unlock
	if err = ks.TimedUnlock(accAddr, passphrase, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if _, err := ks.GetKeyIfUnlockedAndValid(rawToken); err != nil {
		t.Fatal("Account shouldn't have been locked or have any problems returning a valid key")
	}

	// Expiring lock
	time.Sleep(250 * time.Millisecond)
	if _, err := ks.GetKeyIfUnlockedAndValid(rawToken); err != ErrLocked {
		t.Fatal("Account should have been locked")
	}
}

func TestLockAccount(t *testing.T) {
	dir, ks := tmpKeyStore(t)
	defer os.RemoveAll(dir)

	a, err := ks.NewAccount(passphrase)
	if err != nil {
		t.Fatal(err)
	}
	accAddr := a.Address
	var rawToken string
	// We have to issue a token first
	if rawToken, err = ks.IssueTokenForAccount(accAddr, NewTokenClaims("", "")); err != nil {
		t.Fatal(err)
	}

	// Unlock for a period of time
	if err = ks.TimedUnlock(accAddr, passphrase, 5*time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, err := ks.GetKeyIfUnlockedAndValid(rawToken); err != nil {
		t.Fatal("Account shouldn't have been locked or have any problems returning a valid key")
	}
	if err := ks.Lock(accAddr); err != nil {
		t.Fatal(err)
	}
	if _, err := ks.GetKeyIfUnlockedAndValid(rawToken); err != ErrLocked {
		t.Fatal(err)
		t.Fatal("Account should have been locked")
	}
}
