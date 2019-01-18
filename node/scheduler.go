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
	log.Println("start prunning images")
	// TODO: Time has to be a const somewhere
	ticker := time.NewTicker(time.Second * 10)
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

// RemoveImages removes all images
func RemoveImages() {
	summaries, err := manager.GetInstance().ListImages(types.ImageListOptions{All: true})
	common.CheckErr(err, "[RemoveImages] Failed to List images")
	now := time.Now().Unix()
	for _, img := range summaries {
		image := database.ImageLvlDB{}
		imgID := strings.Replace(img.ID, "sha256:", "", -1)

		i, err := database.GetDB().Model(image).Get([]byte(imgID))
		image, ok := i.(database.ImageLvlDB)
		if !ok {
			continue
		}
		// If the image was found into the DB
		if err == nil {
			if time.Unix(image.CreatedTime, 0).Add(common.TenDays).Unix() <= now {
				log.Println("Removing image: ", img.ID)
				manager.GetInstance().RemoveImage(img.ID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
			}
		}
	}
}
