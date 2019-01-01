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

	"github.com/crowdcompute/crowdengine/crypto"
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

// DecryptKey decrypts the private key given a json keystore file
func DecryptKey(password string, data string) (string, error) {

	encjson := encryptedKeyJSON{}
	err := json.Unmarshal([]byte(data), &encjson)
	if err != nil {
		return "", err
	}

	if encjson.Version != ksVersion {
		return "", errors.New("Version Mismatch")
	}
	if encjson.Crypto.Cipher != ksCipher {
		return "", errors.New("Cipher Mismatch")
	}
	mac, err := hex.DecodeString(encjson.Crypto.MAC)
	iv, err := hex.DecodeString(encjson.Crypto.CipherParams.IV)
	salt, err := hex.DecodeString(encjson.Crypto.KDFParams.Salt)
	ciphertext, err := hex.DecodeString(encjson.Crypto.CipherText)
	dk, err := scrypt.Key([]byte(password), salt, encjson.Crypto.KDFParams.N, encjson.Crypto.KDFParams.R, encjson.Crypto.KDFParams.P, encjson.Crypto.KDFParams.DKeyLength)
	hash := crypto.Keccak256(dk[16:32], ciphertext)
	if !bytes.Equal(hash, mac) {
		return "", errors.New("Mac Mismatch")
	}
	aesBlock, err := aes.NewCipher(dk[:16])
	if err != nil {
		return "", err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outputkey := make([]byte, len(ciphertext))
	stream.XORKeyStream(outputkey, ciphertext)
	return hex.EncodeToString(outputkey), nil
}

// EncryptKey encrypts a key using a symmetric algorithm
func EncryptKey(password string, key *Key) (string, error) {
	salt, err := crypto.RandomEntropy(32)
	if err != nil {
		return "", err
	}
	dk, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptKeyLen)
	if err != nil {
		return "", err
	}
	iv, err := crypto.RandomEntropy(aes.BlockSize)
	if err != nil {
		return "", err
	}
	enckey := dk[:16]

	privateKeyBytes, err := hex.DecodeString(key.KeyPair.Private)
	if err != nil {
		return "", err
	}
	aesBlock, err := aes.NewCipher(enckey)
	if err != nil {
		return "", err
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
		Id:      key.Id.String(),
		Version: ksVersion,
	}
	data, err := json.MarshalIndent(&encjson, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
