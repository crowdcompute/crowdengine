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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/crowdcompute/crowdengine/cmd/terminal"
	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/pborman/uuid"
)

// Create generates random keypair
func Create() string {
	keypair, err := crypto.GenerateKeyPair()
	if err != nil {
		log.Fatal(err)
	}

	key := &Key{
		Id:      uuid.NewRandom(),
		KeyPair: &keypair,
	}

	pass, err := terminal.Stdin.GetPassphrase("Please give a password and not forget this password.", true)
	if err != nil {
		log.Fatalf("Error reading passphrase from terminal: %v", err)
	}

	keyDataJSON, err := EncryptKey(pass, key)
	if err != nil {
		log.Fatalf("Error encrypting key: %v", err)
	}

	fileName, err := WriteDataToFile(keyDataJSON, key.KeyPair.Address)
	if err != nil {
		log.Fatalf("Error writing keystore file: %v", err)
	}
	return fileName
}

// WriteDataToFile writes the key data to a file
func WriteDataToFile(data []byte, extra string) (string, error) {
	fileName := createFileName(extra)
	jsonFile, err := os.Create(fileName)
	if err != nil {
		return "nil", err
	}
	defer jsonFile.Close()
	jsonFile.Write(data)
	return fileName, nil
}

// LoadFromFile loads a keystore from file
func LoadFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
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
