package database

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/syndtr/goleveldb/leveldb"
)

// Putter wraps the database write operation supported by both batches and regular databases.
type Putter interface {
	Put(key []byte, value interface{}) error
}

// Deleter wraps the database delete operation supported by both batches and regular databases.
type Deleter interface {
	Delete(key []byte) error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	Putter
	Deleter
	Get(key []byte) (interface{}, error)
	Has(key []byte) (bool, error)
	Close()
}

// DB ...
type DB struct {
	levelDB      *leveldb.DB
	tableName    string
	CurrentModel interface{}
}

var (
	db   *DB
	once sync.Once
)

func GetDB() *DB {
	once.Do(func() {
		lvldb, err := leveldb.OpenFile("database", nil)
		common.CheckErr(err, "[GetDB] Couldn't create a new Level DB")
		db = &DB{levelDB: lvldb}
	})
	return db
}

// Model ...
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

type ImageLvlDB struct {
	Hash        string `json:"hash"`        // The hash of the image
	Signature   string `json:"signature"`   // The uploader of this image
	CreatedTime int64  `json:"createdtime"` // The time the image was created into the current node
}
