package database

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	db   *DB
	once sync.Once
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

// DB represents a level db database
// The keys of the same table are prefixed with the tableName
type DB struct {
	levelDB      *leveldb.DB // The leveldb's handler
	tableName    string      // TableName is the name of the tables created into this db.
	CurrentModel interface{} // The model to store into the db
}

// ImageLvlDB represents the Image Model
type ImageLvlDB struct {
	Hash        string `json:"hash"`        // The hash of the image
	Signature   string `json:"signature"`   // The uploader of this image
	CreatedTime int64  `json:"createdtime"` // The time the image was created into the current node
}
