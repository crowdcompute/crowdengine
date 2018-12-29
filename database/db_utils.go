package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/syndtr/goleveldb/leveldb"
)

// ErrNotFound is returned when no results are returned from the database
var ErrNotFound = errors.New("ErrorNotFound")

// GetDB returns a singleton DB object
func GetDB() *DB {
	once.Do(func() {
		lvldb, err := leveldb.OpenFile("database", nil)
		common.CheckErr(err, "[GetDB] Couldn't create a new Level DB")
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

// Get retrieves a persisted value for a specific key. If there is no results
// ErrNotFound is returned. The provided parameter should be either a byte slice or
// a struct that implements the encoding.BinaryUnmarshaler interface
func (db *DB) Get(key []byte) (interface{}, error) {
	has, err := db.Has(key)
	if err != nil || !has {
		return nil, ErrNotFound
	}

	data, err := db.levelDB.Get(db.prefixKey(key), nil)
	if err == leveldb.ErrNotFound {
		log.Println("image not found in DB")
		return nil, ErrNotFound
	}

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

// Close releases the resources used by the underlying LevelDB.
func (db *DB) Close() {
	db.levelDB.Close()
}
