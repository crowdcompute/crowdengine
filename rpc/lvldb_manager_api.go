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

package rpc

import (
	"context"
	"encoding/json"

	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/log"
)

// LvlDBManagerAPI ...
type LvlDBManagerAPI struct {
}

// NewLvlDBManagerAPI creates a new bootnode API
func NewLvlDBManagerAPI() *LvlDBManagerAPI {
	return &LvlDBManagerAPI{}
}

// GetDBStats returns all values from lvldb
func (api *LvlDBManagerAPI) GetDBStats(ctx context.Context) string {
	stats := database.GetDB().GetProperty("leveldb.stats")
	log.Println(stats)
	return stats
}

func getInstance(objectName string) (interface{}, bool) {
	typeRegistry := make(map[string]interface{})
	typeRegistry["ImageAccount"] = &database.ImageAccount{}
	typeRegistry["Image"] = &database.ImageLvlDB{}
	i, ok := typeRegistry[objectName]
	return i, ok
}

// SelectType returns all values that are of the structName type
func (api *LvlDBManagerAPI) SelectType(ctx context.Context, typeName string) (string, error) {
	var data map[string]string
	var err error
	if i, ok := getInstance(typeName); ok {
		data, err = database.GetDB().Model(i).GetAll()
	}
	if err != nil {
		return "", err
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(dataBytes), err
}

// SelectAll returns all values of the database
func (api *LvlDBManagerAPI) SelectAll(ctx context.Context) (string, error) {
	var data map[string]string
	var err error
	data, err = database.GetDB().GetAll()
	if err != nil {
		return "", err
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(dataBytes), err
}

// SelectImage returns an ImageLvlDB if exists in the database
func (api *LvlDBManagerAPI) SelectImage(ctx context.Context, imgID string) (string, error) {
	image, err := database.GetImageFromDB(imgID)
	if err != nil {
		return "", err
	}
	imgBytes, err := json.Marshal(image)
	log.Println(string(imgBytes))
	return string(imgBytes), nil
}

// SelectImageAccount returns an ImageAccountLvlDB if exists in the database
func (api *LvlDBManagerAPI) SelectImageAccount(ctx context.Context, imgHash string) (string, error) {
	image, err := database.GetImageAccountFromDB(imgHash)
	if err != nil {
		return "", err
	}
	imgBytes, err := json.Marshal(image)
	log.Println(string(imgBytes))
	return string(imgBytes), nil
}
