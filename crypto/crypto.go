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

package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/gogo/protobuf/proto"
	crypto "github.com/libp2p/go-libp2p-crypto"
	gosha3 "golang.org/x/crypto/sha3"
)

// KeyPair represents private/public keys and the public address
type KeyPair struct {
	Private string
	Public  string
	Address string
}

// Sha256Hash hashes data with the sha256
func Sha256Hash(data []byte) hash.Hash {
	d := gosha3.New256()
	d.Write(data)
	return d
}

// HashProtoMsg marshals the proto message and returns a sha256 hash
func HashProtoMsg(message proto.Message) ([]byte, error) {
	bin, err := proto.Marshal(message)
	return Sha256Hash(bin).Sum(nil), err
}

// HashFile hashes the file with the sha256
func HashFile(file *os.File) []byte {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil)
}

// GenerateKeyPair generates a private, public, and address keys
func GenerateKeyPair() (KeyPair, error) {
	priv, pub, err := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	if err != nil {
		return KeyPair{}, err
	}

	privateBytes, err := priv.Bytes()
	if err != nil {
		return KeyPair{}, err
	}

	publicBytes, err := pub.Bytes()
	if err != nil {
		return KeyPair{}, err
	}

	// Drop first x bytes of priv(4), and pub(5)
	privateBytes = privateBytes[4:]
	publicBytes = publicBytes[5:]

	return KeyPair{Private: hex.EncodeToString(privateBytes), Public: hex.EncodeToString(publicBytes), Address: PublicToAddress(publicBytes)}, nil
}

// RestorePrivateKey unmarshals the privateKey
func RestorePrivateKey(privateKey []byte) (crypto.PrivKey, error) {
	return crypto.UnmarshalSecp256k1PrivateKey(privateKey)
}

// RestorePubKey unmarshals the pubKey
func RestorePubKey(pubKey []byte) (crypto.PubKey, error) {
	return crypto.UnmarshalSecp256k1PublicKey(pubKey)
}

// RestorePrivateToKeyPair unmarshals the privateKey and returns a the priv as well the pub keys
func RestorePrivateToKeyPair(privateKey []byte) (crypto.PrivKey, crypto.PubKey, error) {
	priv, err := RestorePrivateKey(privateKey)
	pub := priv.GetPublic()
	if err != nil {
		return priv, pub, err
	}
	return priv, pub, nil
}

// PublicToAddress returns the address of a public key
func PublicToAddress(data []byte) string {
	return hex.EncodeToString(Keccak256(data)[12:])
}

// Keccak256 return sha3 of a given byte array
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
