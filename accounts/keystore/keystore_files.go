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
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/crowdcompute/crowdengine/common"
)

// GetKeystoreFiles returns all path files existent in the keyDir
func GetKeystoreFiles(keyDir string) ([]string, error) {
	var mu sync.RWMutex
	var keyFiles []string
	// List all the files from the keyDir folder
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()

	for _, fi := range files {
		path := filepath.Join(keyDir, fi.Name())
		// Skip any non-key files from the folder
		if nonKeyFile(fi) {
			continue
		}
		// Gather the set of all and fresly modified files
		keyFiles = append(keyFiles, path)
	}
	return keyFiles, nil
}

// nonKeyFile ignores editor backups, hidden files and folders/symlinks.
func nonKeyFile(fi os.FileInfo) bool {
	// Skip editor backups and UNIX-style hidden files.
	if strings.HasSuffix(fi.Name(), "~") || strings.HasPrefix(fi.Name(), ".") {
		return true
	}
	// Skip misc special files, directories (yes, symlinks too).
	if fi.IsDir() || fi.Mode()&os.ModeType != 0 {
		return true
	}
	return false
}

// getAccountsFromPath gets the keystorePath and returns an accounts map if keystore files exist in the path
func getAccountsFromPath(keystorePath string) map[string]Account {
	accounts := make(map[string]Account)
	accountPaths, err := GetKeystoreFiles(keystorePath)
	common.FatalIfErr(err, fmt.Sprintf("There was an error getting keystore files from path: '%s'", keystorePath))
	for _, accountPath := range accountPaths {
		// TODO: Need to decouple. This depends on the accounts filename representation
		accountAddr := strings.Split(accountPath, "--")[2]
		accountAddr = strings.TrimSuffix(accountAddr, path.Ext(accountAddr))
		a := Account{
			Address: accountAddr,
			Path:    accountPath,
		}
		accounts[a.Address] = a
	}
	return accounts
}
