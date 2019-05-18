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

package dockerutil

import(
	"strings"
	"regexp"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/crowdcompute/crowdengine/log"
	"github.com/crowdcompute/crowdengine/manager"
)

// LoadImageToDocker takes a path to an image file and loads it to the docker daemon
func LoadImageToDocker(filePath string) (string, error) {
	log.Println("Loading this image: ", filePath)
	loadImageResp, err := manager.GetInstance().LoadImage(filePath)
	if err != nil {
		return "", err
	}
	log.Println(loadImageResp)
	if imgID, exists := getImageID(loadImageResp); exists {
		// Docker image ID is 64 characters
		return imgID[:64], nil
	}
	// If no image ID exists, we extract the image ID
	// from listing the specific image using its tag
	imageTag := loadImageResp[2 : len(loadImageResp)-5]
	res := getImageSummaryFromTag(imageTag)
	imgID := strings.Replace(res.ID, "sha256:", "", -1)
	log.Println("Loaded image. Image ID: ", imgID)
	return imgID, nil
}

// imageIDExists checks if a docker image ID exists in the loadImageResp.
// Docker image is just after the 'sha256:' prefix
func getImageID(loadImageResp string) (string, bool) {
	r, _ := regexp.Compile("sha256:(.*)")
	matches := r.FindAllStringSubmatch(loadImageResp, -1)
	if len(matches) != 0 {
		return matches[0][1], true
	}
	return "", false
}

// getImageSummaryFromTag returns ImageSummaries from images using a tag
func getImageSummaryFromTag(tag string) types.ImageSummary {
	log.Println(tag)
	fargs := filters.NewArgs()
	fargs.Add("reference", tag)
	res, err := manager.GetInstance().ListImages(
		types.ImageListOptions{
			Filters: fargs,
		})
	if err != nil {
		log.Println("error: ", err)
	}
	return res[0] // we know that docker tag is unique thus returning only one summary
}
