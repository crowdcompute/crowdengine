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

// GetImageAccountFromDB returns an ImageAccount if exists in the database
func GetImageAccountFromDB(hash string) (*ImageAccount, error) {
	image := &ImageAccount{}

	i, err := GetDB().Model(image).Get([]byte(hash))
	if err != nil {
		return nil, err
	}
	image = i.(*ImageAccount)
	return image, nil
}

// GetImageFromDB returns an ImageLvlDB if exists in the database
func GetImageFromDB(imgHash string) (*ImageLvlDB, error) {
	image := &ImageLvlDB{}
	i, err := GetDB().Model(image).Get([]byte(imgHash))
	if err != nil {
		return nil, err
	}
	image = i.(*ImageLvlDB)
	return image, nil
}
