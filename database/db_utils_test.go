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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	imagePut := ImageLvlDB{Hash: "test", Signature: "test", CreatedTime: time.Now().Unix()}
	GetDB().Model(&imagePut).Put([]byte("testKey"))
}

func TestPutGet(t *testing.T) {
	imageGet := ImageLvlDB{}
	i, err := GetDB().Model(&imageGet).Get([]byte("testKey"))
	if err == nil {
		imageGet, ok := i.(*ImageLvlDB)
		if !ok {
			t.Errorf("Type assertion error")
		}
		assert.Equal(t, imageGet.Hash, "test")
		assert.Equal(t, imageGet.Signature, "test")
	} else {
		t.Errorf("Couldn't get the image")
	}
}

func TestHas(t *testing.T) {
	has, err := GetDB().Model(&ImageLvlDB{}).Has([]byte("testKey"))
	if err != nil {
		t.Errorf("Error getting image")
	}
	assert.True(t, has)
}

func TestDelete(t *testing.T) {
	err := GetDB().Model(&ImageLvlDB{}).Delete([]byte("testKey"))
	imageGet := ImageLvlDB{}
	_, err = GetDB().Model(&imageGet).Get([]byte("testKey"))
	if err == nil {
		t.Errorf("Got a deleted image")
	}
}
