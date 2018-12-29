package database

import (
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	image := ImageLvlDB{Hash: "hash", Signature: "signature", CreatedTime: time.Now().Unix()}
	GetDB().Model(&image).Put([]byte("imageID"))

	imageGet := ImageLvlDB{}
	i, err := GetDB().Model(&imageGet).Get([]byte("imageID"))
	if err == nil {
		imageGet, ok := i.(*ImageLvlDB)
		if !ok {
			return
		}
		fmt.Println(imageGet)
	} else {
		fmt.Println(err)
	}
}
