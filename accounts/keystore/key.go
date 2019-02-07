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
	"log"
	"path/filepath"

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/pborman/uuid"
)

// Key represents a key with UUID
type Key struct {
	*crypto.KeyPair
	ID uuid.UUID
}

// NewKey creates a new Key
func NewKey() *Key {
	keypair, err := crypto.GenerateKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	return &Key{
		ID:      uuid.NewRandom(),
		KeyPair: &keypair,
	}
}

// NewKeyAndStoreToFile creates a new Key
func NewKeyAndStoreToFile(passphrase string) (*Key, string) {
	key := NewKey()
	return key, key.StoreKeyToFile(passphrase)
}

// StoreKeyToFile generates random keypair
func (key *Key) StoreKeyToFile(pass string) string {
	keyDataJSON, err := MarshalKey(pass, key)
	if err != nil {
		log.Fatalf("Error encrypting key: %v", err)
	}
	cfg := config.DefaultConfig()
	fileName, err := common.WriteDataToFile(keyDataJSON, filepath.Join(cfg.Global.KeystoreDir, createFileName(key.KeyPair.Address)))
	if err != nil {
		log.Fatalf("Error writing keystore file: %v", err)
	}
	return fileName
}
