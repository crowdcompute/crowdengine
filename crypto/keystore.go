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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/crowdcompute/crowdengine/cmd/terminal"
	"github.com/pborman/uuid"
)

// Create generates random keypair
func Create() {
	keypair, err := GenerateKeyPair()
	if err != nil {
		log.Fatal(err)
	}

	key := &Key{
		Id:      uuid.NewRandom(),
		KeyPair: &keypair,
	}

	pass, err := terminal.Stdin.GetPassphrase("Passphrase: ", true)
	if err != nil {
		log.Fatalf("Error reading passphrase from terminal: %v", err)
	}

	keyJSON, err := EncryptKey(pass, key)
	if err != nil {
		log.Fatalf("Error encrypting key: %v", err)
	}

	err = WriteToFile(keyJSON, key.KeyPair.Address)
	if err != nil {
		log.Fatalf("Error writing keystore file: %v", err)
	}

}

// WriteToFile writes to a file
func WriteToFile(keydata string, address string) error {

	jsonFile, err := os.Create(keyFileName(address))
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	jsonFile.Write([]byte(keydata))
	return nil
}

// LoadFromFile loads a keystore from file
func LoadFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func keyFileName(address string) string {
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
