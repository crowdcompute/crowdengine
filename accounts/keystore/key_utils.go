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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/scrypt"
)

var (
	nameKDF      = "scrypt"
	scryptKeyLen = 32
	scryptN      = 1 << 18
	scryptR      = 8
	scryptP      = 1
	ksVersion    = 3
	ksCipher     = "aes-128-ctr"
)

// UnmarshalKey decrypts the private key given a passphrase and a json keystore file
func UnmarshalKey(data []byte, passphrase string) (*Key, error) {
	encjson := encryptedKeyJSON{}
	err := json.Unmarshal(data, &encjson)
	if err != nil {
		return &Key{}, err
	}
	if encjson.Version != ksVersion {
		return &Key{}, errors.New("Version Mismatch")
	}
	if encjson.Crypto.Cipher != ksCipher {
		return &Key{}, errors.New("Cipher Mismatch")
	}
	mac, err := hex.DecodeString(encjson.Crypto.MAC)
	iv, err := hex.DecodeString(encjson.Crypto.CipherParams.IV)
	salt, err := hex.DecodeString(encjson.Crypto.KDFParams.Salt)
	ciphertext, err := hex.DecodeString(encjson.Crypto.CipherText)
	dk, err := scrypt.Key([]byte(passphrase), salt, encjson.Crypto.KDFParams.N, encjson.Crypto.KDFParams.R, encjson.Crypto.KDFParams.P, encjson.Crypto.KDFParams.DKeyLength)
	hash := crypto.Keccak256(dk[16:32], ciphertext)
	if !bytes.Equal(hash, mac) {
		return &Key{}, errors.New("Mac Mismatch")
	}
	aesBlock, err := aes.NewCipher(dk[:16])
	if err != nil {
		return &Key{}, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outputkey := make([]byte, len(ciphertext))
	stream.XORKeyStream(outputkey, ciphertext)
	privKey, err := crypto.RestorePrivateKey(outputkey)

	return &Key{
		ID: uuid.UUID(encjson.ID),
		KeyPair: &crypto.KeyPair{
			Private: privKey,
			Address: encjson.Address,
		},
	}, nil
}

// MarshalKey encrypts a key using a symmetric algorithm
func MarshalKey(passphrase string, key *Key) ([]byte, error) {
	salt, err := crypto.RandomEntropy(32)
	if err != nil {
		return nil, err
	}
	dk, err := scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, scryptKeyLen)
	if err != nil {
		return nil, err
	}
	iv, err := crypto.RandomEntropy(aes.BlockSize)
	if err != nil {
		return nil, err
	}
	enckey := dk[:16]

	privateKeyBytes, err := key.KeyPair.Private.Bytes()
	privateKeyBytes = privateKeyBytes[4:]
	if err != nil {
		return nil, err
	}
	aesBlock, err := aes.NewCipher(enckey)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	cipherText := make([]byte, len(privateKeyBytes))
	stream.XORKeyStream(cipherText, privateKeyBytes)

	mac := crypto.Keccak256(dk[16:32], cipherText)
	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	sp := ScryptParams{
		N:          scryptN,
		R:          scryptR,
		P:          scryptP,
		DKeyLength: scryptKeyLen,
		Salt:       hex.EncodeToString(salt),
	}

	keyjson := cryptoJSON{
		Cipher:       ksCipher,
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          nameKDF,
		KDFParams:    sp,
		MAC:          hex.EncodeToString(mac),
	}

	encjson := encryptedKeyJSON{
		Address: key.KeyPair.Address,
		Crypto:  keyjson,
		ID:      key.ID.String(),
		Version: ksVersion,
	}
	data, err := json.MarshalIndent(&encjson, "", "  ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Returns a name joining the timestamp and the address
func createFileName(address string) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s.json", toISO8601(ts), address)
}

func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}
