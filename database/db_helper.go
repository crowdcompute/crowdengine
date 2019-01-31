package database

import (
	"fmt"
)

// GetImageAccountFromDB returns an ImageAccount if exists in the database
func GetImageAccountFromDB(hash string) (*ImageAccount, error) {
	image := &ImageAccount{}
	i, err := GetDB().Model(image).Get([]byte(hash))
	if err != nil && err != ErrNotFound {
		return nil, fmt.Errorf("There was an error getting the image from lvldb")
	}
	image = i.(*ImageAccount)
	return image, nil
}

// GetImageFromDB returns an ImageLvlDB if exists in the database
func GetImageFromDB(imgID string) (*ImageLvlDB, error) {
	image := &ImageLvlDB{}
	i, err := GetDB().Model(image).Get([]byte(imgID))
	if err != nil && err != ErrNotFound {
		return nil, fmt.Errorf("There was an error getting the image from lvldb")
	}
	image = i.(*ImageLvlDB)
	return image, nil
}
