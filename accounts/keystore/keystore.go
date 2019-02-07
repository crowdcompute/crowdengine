package keystore

import (
	"fmt"
	"sync"
	"time"

	"github.com/crowdcompute/crowdengine/common"
	crypto "github.com/libp2p/go-libp2p-crypto"
)

// Account represents a Crowd Compute account located at a specific location Path
type Account struct {
	Address string `json:"address"` // Crowd Compute account address derived from the key
	Path    string `json:"path"`    // Optional resource locator within a backend
}

type KeyStore struct {
	accounts map[string][]Account // address -> Account
	unlocked map[string]*unlocked // Currently unlocked account (decrypted private keys)

	mu sync.RWMutex
}

type unlocked struct {
	*Key
	abort chan struct{}
}

func NewKeyStore() *KeyStore {
	return &KeyStore{
		accounts: make(map[string][]Account),
	}
}

// NewAccount generates a new key and stores it into the key directory,
// encrypting it with the passphrase.
func (ks *KeyStore) NewAccount(passphrase string) (Account, error) {
	key, fileName := NewKeyAndStoreToFile(passphrase)
	a := Account{
		Address: key.Address,
		Path:    fileName,
	}
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.accounts[a.Address] = append(ks.accounts[a.Address], a)
	return a, nil
}

// Accounts returns all key files present in the directory.
func (ks *KeyStore) Accounts(address string) []Account {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	cpy := make([]Account, len(ks.accounts[address]))
	copy(cpy, ks.accounts[address])
	return cpy
}

// Lock removes the private key with the given address from memory.
func (ks *KeyStore) Lock(addr string) error {
	ks.mu.Lock()
	if unl, found := ks.unlocked[addr]; found {
		ks.mu.Unlock()
		ks.expire(addr, unl, time.Duration(0)*time.Nanosecond)
	} else {
		ks.mu.Unlock()
	}
	return nil
}

// Unlock unlocks the given account indefinitely.
func (ks *KeyStore) Unlock(a Account, passphrase string) error {
	return ks.TimedUnlock(a, passphrase, 0)
}

// TimedUnlock unlocks the given account with the passphrase. The account
// stays unlocked for the duration of timeout. A timeout of 0 unlocks the account
// until the program exits. The account must match a unique key file.
//
// If the account address is already unlocked for a duration, TimedUnlock extends or
// shortens the active unlock timeout. If the address was previously unlocked
// indefinitely the timeout is not altered.
func (ks *KeyStore) TimedUnlock(a Account, passphrase string, timeout time.Duration) error {
	a, key, err := ks.getDecryptedKey(a, passphrase)
	if err != nil {
		return err
	}

	ks.mu.Lock()
	defer ks.mu.Unlock()
	u, found := ks.unlocked[a.Address]
	if found {
		if u.abort == nil {
			// The address was unlocked indefinitely, so unlocking
			// it with a timeout would be confusing.
			zeroKey(key.Private)
			return nil
		}
		// Terminate the expire goroutine and replace it below.
		close(u.abort)
	}
	if timeout > 0 {
		u = &unlocked{Key: key, abort: make(chan struct{})}
		go ks.expire(a.Address, u, timeout)
	} else {
		u = &unlocked{Key: key}
	}
	ks.unlocked[a.Address] = u
	return nil
}

func (ks *KeyStore) expire(addr string, u *unlocked, timeout time.Duration) {
	t := time.NewTimer(timeout)
	defer t.Stop()
	select {
	case <-u.abort:
		// just quit
	case <-t.C:
		ks.mu.Lock()
		// only drop if it's still the same key instance that dropLater
		// was launched with. we can check that using pointer equality
		// because the map stores a new pointer every time the key is
		// unlocked.
		if ks.unlocked[addr] == u {
			zeroKey(u.Private)
			delete(ks.unlocked, addr)
		}
		ks.mu.Unlock()
	}
}

func (ks *KeyStore) getDecryptedKey(a Account, pass string) (Account, *Key, error) {
	a, err := ks.Find(a)
	if err != nil {
		return a, nil, err
	}
	key, err := ks.GetKey(a.Address, a.Path, pass)
	return a, key, err
}

// Find resolves the given account into a unique entry in the keystore.
func (ks *KeyStore) Find(a Account) (Account, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	addressAccounts := ks.accounts[a.Address]
	for _, acc := range addressAccounts {
		if acc.Path == a.Path {
			return acc, nil
		}
	}
	return a, fmt.Errorf("Could not find address")
}

// GetKey loads the key from the keystore and decrypt its contents
func (ks *KeyStore) GetKey(addr string, filename, pass string) (*Key, error) {
	keyjson, err := common.LoadFromFile(filename)
	if err != nil {
		return nil, err
	}
	key, err := UnmarshalKey(keyjson, pass)
	if err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, addr)
	}
	return key, nil
}

// zeroKey zeroes a private key in memory.
func zeroKey(k crypto.PrivKey) {
	// func zeroKey(k string) {
	// k = nil
	// b := k.D.Bits()
	// for i := range b {
	// 	b[i] = 0
	// }
}
