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

package common

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var r *rand.Rand // Rand for this package.

// check if applicable or should be placed inside random string
func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString generates a random string of strlen length
func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""
	for i := 0; i < strlen; i++ {
		index := r.Intn(len(chars))
		result += chars[index : index+1]
	}
	return result
}

// FillString fills the value string with toFill string for up to upToLength
func FillString(value string, upToLength int) string {
	for {
		currLen := len(value)
		if currLen < upToLength {
			value = value + FillChar
			continue
		}
		break
	}
	return value
}

// RemoveFile removes the filePath file from the os
func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("There was an error removing the file %s. Error: %s", filePath, err)
	}
	return nil
}
