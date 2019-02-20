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

package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// ErrNotFound is returned when no results are returned from the database
var ErrNotFound = errors.New("ErrorNotFound")
var lvlDBPath = "gocc_db"

// SetLvlDBPath sets the path for lvldb to store the data
func SetLvlDBPath(path string) {
	lvlDBPath = filepath.Join(path, lvlDBPath)
}

// GetDB returns a singleton DB object
func GetDB() *DB {
	once.Do(func() {
		lvldb, err := leveldb.OpenFile(lvlDBPath, nil)
		common.FatalIfErr(err, "Couldn't create a new Level DB")
		db = &DB{levelDB: lvldb}
	})
	return db
}

// Model gets an  interface iface and returns a creates a DB object.
// The tableName and CurrentModel are extracted of the iface
func (db *DB) Model(iface interface{}) *DB {
	return &DB{
		tableName:    fmt.Sprintf("%T", iface),
		CurrentModel: iface,
		levelDB:      db.levelDB,
	}
}

// GetAll returns all values in the database if no db.tableName is given,
// else it returns those values that are of the db.tableName type
func (db *DB) GetAll() (map[string]string, error) {
	var iter iterator.Iterator
	if db.tableName == "" {
		iter = db.levelDB.NewIterator(nil, nil)
	} else {
		iter = db.levelDB.NewIterator(util.BytesPrefix([]byte(db.tableName)), nil)
	}
	data := make(map[string]string)
	for iter.Next() {
		key := string(iter.Key())
		if strings.HasPrefix(key, db.tableName) {
			value := iter.Value()
			log.Println(string(value))
			data[key] = string(value)
		}
	}
	iter.Release()
	err := iter.Error()
	return data, err
}

// Get retrieves a persisted value for a specific key. If there is no results
// ErrNotFound is returned. The provided parameter should be either a byte slice or
// a struct that implements the encoding.BinaryUnmarshaler interface
func (db *DB) Get(key []byte) (interface{}, error) {
	has, err := db.Has(key)
	if err != nil || !has {
		return nil, ErrNotFound
	}
	data, _ := db.levelDB.Get(db.prefixKey(key), nil)
	err = json.Unmarshal(data, db.CurrentModel)
	return db.CurrentModel, err
}

// Has returns true if a key exists in the lvldb database
func (db *DB) Has(key []byte) (bool, error) {
	return db.levelDB.Has(db.prefixKey(key), nil)
}

// Put stores an object that implements Binary for a specific key.
func (db *DB) Put(key []byte) (err error) {
	bytes := []byte{}
	if bytes, err = json.Marshal(db.CurrentModel); err != nil {
		return err
	}

	return db.levelDB.Put(db.prefixKey(key), bytes, nil)
}

// prefixKey prefixes all the key with the tableName
func (db *DB) prefixKey(key []byte) []byte {
	return append([]byte(db.tableName), key...)
}

// Delete removes entries stored under a specific key.
func (db *DB) Delete(key []byte) error {
	return db.levelDB.Delete(db.prefixKey(key), nil)
}

// GetProperty returns value of the given property name.
func (db *DB) GetProperty(name string) string {
	res, err := db.levelDB.GetProperty(name)
	common.FatalIfErr(err, "Failed to read database stats")
	return res
}

// Close releases the resources used by the underlying LevelDB.
func (db *DB) Close() {
	db.levelDB.Close()
}
