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

// ImageAccount represents the Image Account Model
// These are the images uploaded to the dev node by users
// TODO: name to be changed
type ImageAccount struct {
	Signature   string `json:"signature"`   // The uploader of this image
	Path        string `json:"path"`        // Physical path of the location of the image
	CreatedTime int64  `json:"createdtime"` // The time the image was created into the current node
}
