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
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
	"unicode"
)

var r *rand.Rand // Rand for this package.

// check if applicable or should be placed inside random string
func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// LcFirst converts the first letter of s to lowercase
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
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

// FileExist checks if a file exists at filePath.
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// WriteDataToFile writes the data to a file named fileName
func WriteDataToFile(data []byte, filePath string) (string, error) {
	// Create the directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0600
	if err := os.MkdirAll(filepath.Dir(filePath), dirPerm); err != nil {
		return "", err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return "nil", err
	}
	defer file.Close()
	file.Write(data)
	return filePath, nil
}

// LoadDataFromFile loads data from file
func LoadDataFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}
