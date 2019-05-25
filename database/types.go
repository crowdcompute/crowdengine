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

// ImageLoadDocker represents the Image Model. Keeps track of the images loaded into docker
// Usage: Whenever we load an image to docker we store it in our DB as well
//    	  We remove all images that are in the node for a long period of time
//		  We keep track who loaded the image to the node's docker engine and use
//		  their signature (of the Hash) to identify them as the owners of their images
type ImageLoadDocker struct {
	Hash        string   `json:"hash"`        // The hash of the image
	Signatures  []string `json:"signatures"`  // Signature Verifies the uploader of this image. Same image might have multiple uploaders
	CreatedTime int64    `json:"createdtime"` // The time the image was loaded into the current node's docker engine
}

// ImageAccount represents the Image Account Model. Keeps track of the files uploaded via the Fileserver
// Usage: Dev nodes store this information about the user who uploaded the image
// TODO: name to be changed
type ImageAccount struct {
	Signature   string `json:"signature"`   // The uploader of this image
	Path        string `json:"path"`        // Physical path of the location of the image
	CreatedTime int64  `json:"createdtime"` // The time the image was loaded into the current node's docker engine
}
