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

package node

import (
	"strings"
	"time"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/log"
	"github.com/crowdcompute/crowdengine/manager"
	"github.com/docker/docker/api/types"
)

// PruneImages checks if there are any images to be removed based on a time interval
// Running for ever, or until node dies
func PruneImages(quit <-chan struct{}) {
	ticker := time.NewTicker(common.RemoveImagesInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Checking if there are images to be removed...")
			RemoveImages()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

// RemoveImages removes all images that got expired
// This is a goroutine
func RemoveImages() {
	summaries, err := manager.GetInstance().ListImages(types.ImageListOptions{All: true})
	if err != nil {
		log.Println("There is an error listing images. Stopped checking for expired images... Error : ", err)
		return
	}
	for _, imgSummary := range summaries {
		imgID := extractImgID(imgSummary)
		if image, ok := getImageFromDB(imgID); ok {
			if imageExpired(image.CreatedTime) {
				log.Println("Removing image with ID: ", imgID)
				removeImageFromDocker(imgID)
				removeImageFromDB(imgID)
			}
		}
	}
}

func extractImgID(imgSummary types.ImageSummary) string {
	return strings.Replace(imgSummary.ID, "sha256:", "", -1)
}

func getImageFromDB(imgID string) (database.ImageLvlDB, bool) {
	image := database.ImageLvlDB{}
	i, err := database.GetDB().Model(image).Get([]byte(imgID))
	common.FatalIfErr(err, "There was an error getting the image from lvldb")
	image, ok := i.(database.ImageLvlDB)
	return image, ok
}

func imageExpired(createdTime int64) bool {
	now := time.Now().Unix()
	return time.Unix(createdTime, 0).Add(common.TenDays).Unix() <= now
}

func removeImageFromDocker(imgID string) {
	_, err := manager.GetInstance().RemoveImage(imgID,
		types.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		})
	common.FatalIfErr(err, "There was an error removing the image from docker")
}

func removeImageFromDB(imgID string) {
	err := database.GetDB().Model(database.ImageLvlDB{}).Delete([]byte(imgID))
	common.FatalIfErr(err, "There was an error deleting the image from lvldb")
}
