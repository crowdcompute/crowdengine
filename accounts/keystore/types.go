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

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

// ScryptParams is used to record the params of KDF
type ScryptParams struct {
	N          int    `json:"n"`
	R          int    `json:"r"`
	P          int    `json:"p"`
	DKeyLength int    `json:"dklen"`
	Salt       string `json:"salt"`
}

type cryptoJSON struct {
	Cipher       string           `json:"cipher"`
	CipherText   string           `json:"ciphertext"`
	CipherParams cipherparamsJSON `json:"cipherparams"`
	KDF          string           `json:"kdf"`
	KDFParams    ScryptParams     `json:"kdfparams"`
	MAC          string           `json:"mac"`
}

type encryptedKeyJSON struct {
	Address string     `json:"address"`
	Crypto  cryptoJSON `json:"crypto"`
	Id      string     `json:"id"`
	Version int        `json:"version"`
}
