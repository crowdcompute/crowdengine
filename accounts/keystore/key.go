package keystore

import (
	crypto "github.com/libp2p/go-libp2p-crypto"
)

type Key struct {
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	PrivateKey crypto.PrivKey
}
