// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"encoding/json"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

// ErrNotFound is returned when no results are returned from the database
var ErrNotFound = errors.New("ErrorNotFound")

// ErrInvalidArgument is returned when the argument type does not match the expected type
var ErrInvalidArgument = errors.New("ErrorInvalidArgument")

// DBStore uses LevelDB to store values.
type DBStore struct {
	db *leveldb.DB
}

// NewDBStore creates a new instance of DBStore.
func NewDBStore(path string) (s *DBStore, err error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DBStore{
		db: db,
	}, nil
}

// Get retrieves a persisted value for a specific key. If there is no results
// ErrNotFound is returned. The provided parameter should be either a byte slice or
// a struct that implements the encoding.BinaryUnmarshaler interface
func (s *DBStore) Get(key []byte) (i interface{}, err error) {
	has, err := s.db.Has(key, nil)
	if err != nil || !has {
		return nil, ErrNotFound
	}

	data, err := s.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, ErrNotFound
	}
	err = json.Unmarshal(data, i)
	return
}

func (s *DBStore) Has(key []byte) (bool, error) {
	return s.db.Has(key, nil)
}

// Put stores an object that implements Binary for a specific key.
func (s *DBStore) Put(key []byte, i interface{}) (err error) {
	bytes := []byte{}

	if bytes, err = json.Marshal(i); err != nil {
		return err
	}

	return s.db.Put(key, bytes, nil)
}

// Delete removes entries stored under a specific key.
func (s *DBStore) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

// Close releases the resources used by the underlying LevelDB.
func (s *DBStore) Close() {
	s.db.Close()
}
