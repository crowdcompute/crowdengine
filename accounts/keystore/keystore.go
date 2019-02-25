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
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/crowdcompute/crowdengine/log"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrLocked              = errors.New("Account is locked")
	ErrCouldNotFindAccount = errors.New("Could not find account")
	ErrTokenNotIssued      = errors.New("Token not issued for this account yet")
)

// Account represents a Crowd Compute account located at a specific location Path
type Account struct {
	Token   *jwt.Token
	Address string `json:"address"` // Crowd Compute account address derived from the key
	Path    string `json:"path"`    // Optional resource locator within a backend
}

type KeyStore struct {
	keyDir      string
	accounts    map[string]Account   // address -> Account
	unlockedAcc map[string]*unlocked // Currently unlocked accounts
	symmKey     []byte               // Key for signing tokens
	mu          sync.RWMutex
}

type unlocked struct {
	*Key
	abort chan struct{}
}

// NewKeyStore creates and returns a new keystore
func NewKeyStore(keyDir string) *KeyStore {
	// TODO: give the appropriate permissions here
	const dirPerm = 0777
	if err := os.MkdirAll(keyDir, dirPerm); err != nil {
		return nil
	}
	symmKey, err := crypto.RandomEntropy(32)
	common.FatalIfErr(err, "There was an error getting random entropy")
	return &KeyStore{
		accounts:    make(map[string]Account),
		unlockedAcc: make(map[string]*unlocked),
		symmKey:     symmKey,
		keyDir:      keyDir,
	}
}

// NewAccount generates a new key and stores it into the key directory,
// encrypting it with the passphrase.
func (ks *KeyStore) NewAccount(passphrase string) (Account, error) {
	key, fileName := NewKeyAndStoreToFile(passphrase, ks.keyDir)
	a := Account{
		Address: key.Address,
		Path:    fileName,
	}
	ks.addAccount(a)
	return a, nil
}

// addAccount adds or replaces an account
func (ks *KeyStore) addAccount(a Account) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.accounts[a.Address] = a
}

// Delete deletes an account from the disk and the memory
func (ks *KeyStore) Delete(address, passphrase string) error {
	var err error
	a, err := ks.Find(address)
	if err != nil {
		return err
	}
	// Decrypting the key isn't really necessary, but we do
	// it anyway to check the passphrase
	if _, err = ks.extractKeyFromFile(a.Address, a.Path, passphrase); err != nil {
		return err
	}
	err = os.Remove(a.Path)
	if err == nil {
		ks.deleteAccount(address)
	}
	return err
}

// deleteAccount removes an account from the accounts map
func (ks *KeyStore) deleteAccount(address string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	delete(ks.accounts, address)
}

// GetAccounts returns all accounts in this keystore
func (ks *KeyStore) GetAccounts() []Account {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	accounts := make([]Account, 0)
	for _, acc := range ks.accounts {
		accounts = append(accounts, acc)
	}
	return accounts
}

// IssueTokenForAccount issues a token for the specified account
func (ks *KeyStore) IssueTokenForAccount(accAddress string, tClaims *TokenClaims) (string, error) {
	a, err := ks.Find(accAddress)
	if err != nil {
		return "", err
	}
	tok, err := NewToken(ks.symmKey, tClaims)
	a.Token = tok
	ks.addAccount(a)
	return tok.Raw, err
}

// GetKeyIfUnlockedAndValid returns the Key structure of the account if the account is unlocked
// and its token is valid
func (ks *KeyStore) GetKeyIfUnlockedAndValid(rawToken string) (*Key, error) {
	if verified, err := VerifyToken(rawToken, ks.symmKey); !verified {
		return nil, fmt.Errorf("Couldn't verify token: {%s}", rawToken)
	} else if err != nil {
		return nil, err
	}
	ks.mu.Lock()
	defer ks.mu.Unlock()
	unlockedKey, found := ks.unlockedAcc[HashToken(rawToken)]
	if !found {
		return nil, ErrLocked
	}
	return unlockedKey.Key, nil
}

// Find resolves the given account into a unique entry in the keystore.
func (ks *KeyStore) Find(accAddress string) (Account, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	if acc, found := ks.accounts[accAddress]; found {
		return acc, nil
	}
	log.Errorf("%s : {%s}", ErrCouldNotFindAccount.Error(), accAddress)
	return Account{}, ErrCouldNotFindAccount
}

// Lock removes the private key with the given address from memory.
func (ks *KeyStore) Lock(accAddress string) error {
	a, err := ks.Find(accAddress)
	if err != nil {
		return err
	}
	hashedToken := HashToken(a.Token.Raw)
	ks.mu.Lock()
	if unl, found := ks.unlockedAcc[hashedToken]; found {
		// Terminate the routine
		if unl.abort != nil {
			close(unl.abort)
		}
		delete(ks.unlockedAcc, hashedToken)
		ks.mu.Unlock()
	} else {
		ks.mu.Unlock()
	}
	return nil
}

// Unlock unlocks the given account indefinitely.
func (ks *KeyStore) Unlock(accAddress, passphrase string) error {
	return ks.TimedUnlock(accAddress, passphrase, 0)
}

// TimedUnlock unlocks the given account with the passphrase. The account
// stays unlocked for the duration of timeout. A timeout of 0 unlocks the account
// until the program exits. The account must match a unique key file.
//
// If the account address is already unlocked for a duration, TimedUnlock extends or
// shortens the active unlock timeout. If the address was previously unlocked
// indefinitely the timeout is not altered.
func (ks *KeyStore) TimedUnlock(accAddress, passphrase string, timeout time.Duration) error {
	a, err := ks.Find(accAddress)
	if err != nil {
		return err
	}
	key, err := ks.extractKeyFromFile(a.Address, a.Path, passphrase)
	if err != nil {
		return err
	}
	hashedToken := HashToken(a.Token.Raw)

	ks.mu.Lock()
	defer ks.mu.Unlock()
	u, found := ks.unlockedAcc[hashedToken]
	if found {
		if u.abort == nil {
			// The address was unlocked indefinitely, so unlocking
			// it with a timeout would be confusing.
			return nil
		}
		// Terminate the expire goroutine and replace it below.
		close(u.abort)
	}
	if timeout > 0 {
		u = &unlocked{Key: key, abort: make(chan struct{})}
		go ks.expire(hashedToken, u, timeout)
	} else {
		u = &unlocked{Key: key}
	}
	ks.unlockedAcc[hashedToken] = u
	return nil
}

func (ks *KeyStore) expire(hashedToken string, u *unlocked, timeout time.Duration) {
	t := time.NewTimer(timeout)
	defer t.Stop()
	select {
	case <-u.abort:
		// just quit
	case <-t.C:
		log.Printf("The account has expired. Locking...")
		ks.mu.Lock()
		// only drop if it's still the same key instance that dropLater
		// was launched with. we can check that using pointer equality
		// because the map stores a new pointer every time the key is
		// unlocked.
		if ks.unlockedAcc[hashedToken] == u {
			delete(ks.unlockedAcc, hashedToken)
		}
		ks.mu.Unlock()
	}
}

// extractKeyFromFile loads the key from the account's path and unmarshals its contents
func (ks *KeyStore) extractKeyFromFile(accAddress, accPath, pass string) (*Key, error) {
	// gets the byte data from a file
	keyData, err := common.LoadDataFromFile(accPath)
	if err != nil {
		return nil, err
	}
	// Unmarshals the json data to a struct
	key, err := UnmarshalKey(keyData, pass)
	if err != nil {
		return nil, fmt.Errorf("Given passphrase could be wrong. Error: %s", err)
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != accAddress {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, accAddress)
	}
	return key, nil
}
