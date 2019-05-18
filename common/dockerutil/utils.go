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
	"encoding/json"
	"encoding/hex"
	"strings"
	"regexp"
	"fmt"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/crowdcompute/crowdengine/log"
	"github.com/crowdcompute/crowdengine/manager"
	"github.com/crowdcompute/crowdengine/database"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/crowdcompute/crowdengine/crypto"
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

// InspectContainerRaw inspects a container returning raw results
func InspectContainerRaw(containerID string) ([]byte, error) {
	log.Println("Inspecting this container: ", containerID)
	getSize := true
	inspection, rawData, err := manager.GetInstance().InspectContainerRaw(containerID, getSize)
	log.Printf("Result inspection the container %t\n", inspection.State.Running)
	return rawData, err
}

// GetRawContainersForUser returns a raw list of containers for the user with the specific publicKey
func GetRawContainersForUser(publicKey string) (string, error){
	containers, err := ListContainersForUser(publicKey)
	if err != nil {
		log.Println(err, "Error listing containers for user.")
		return "", err
	}
	containersBytes, err := json.Marshal(containers)
	log.Println("Container summaries:", string(containersBytes))
	return string(containersBytes), err
}

// ListContainersForUser returns a list of containers for the user with the specific publicKey
func ListContainersForUser(publicKey string) ([]types.Container, error) {
	containers := make([]types.Container, 0)
	allContainers, err := manager.GetInstance().ListContainers()
	if err != nil {
		return nil, fmt.Errorf("Error listing images. Error: %v", err)
	}

	for _, container := range allContainers {
		hash, signatures, err := getImgDataFromDB(container.ImageID)
		if err != nil {
			if err == database.ErrNotFound {
				log.Println("Continuing... ")
				continue
			}
			return nil, err
		}
		// Verify all signatures for the same image
		for _, signature := range signatures {
			signedBytes, err := hex.DecodeString(signature)
			if err != nil {
				return nil, err
			}
			if ok, err := verifyUser(publicKey, hash, signedBytes); ok && err == nil {
				containers = append(containers, container)
				// TODO: Delete those comments. Only for debugging mode
				// } else if !ok {
				// 	log.Println("Could not verify this user. Signature could not be verified by the Public key...")
			} else if err != nil {
				return nil, err
			}
		}
	}
	return containers, nil
}

// GetRawImagesForUser a raw list of images for the user with the specific publicKey
func GetRawImagesForUser(publicKey string) (string, error){
	images, err := ListImagesForUser(publicKey)
	if err != nil {
		log.Println(err, "Error listing containers for user.")
		return "", err
	}
	imagesBytes, err := json.Marshal(images)
	log.Println("Images summaries:", string(imagesBytes))
	return string(imagesBytes), err
}

// ListImagesForUser list images for the user with the specific publicKey
func ListImagesForUser(publicKey string) ([]types.ImageSummary, error) {
	imgSummaries := make([]types.ImageSummary, 0)
	allSummaries, err := manager.GetInstance().ListImages(types.ImageListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("Error listing images. Error: %v", err)
	}

	for _, imgSummary := range allSummaries {
		hash, signatures, err := getImgDataFromDB(imgSummary.ID)
		if err != nil {
			if err == database.ErrNotFound {
				log.Println("Continuing... ")
				continue
			}
			return nil, err
		}
		// Verify all signatures for the same image
		for _, signature := range signatures {
			signedBytes, err := hex.DecodeString(signature)
			if err != nil {
				return nil, err
			}
			if ok, err := verifyUser(publicKey, hash, signedBytes); ok && err == nil {
				imgSummaries = append(imgSummaries, imgSummary)
				// TODO: Delete those comments. Only for debugging mode
				// } else if !ok {
				// 	log.Println("Could not verify this user. Signature could not be verified by the Public key...")
			} else if err != nil {
				return nil, err
			}
		}
	}
	return imgSummaries, nil
}

func getImgDataFromDB(imgID string) ([]byte, []string, error) {
	imgID = strings.Replace(imgID, "sha256:", "", -1)
	if image, err := database.GetImageFromDB(imgID); err == nil {
		hashBytes, err := hex.DecodeString(image.Hash)
		if err != nil {
			return nil, nil, err
		}
		return hashBytes, image.Signatures, err
	} else {
		return nil, nil, err
	}
}

func verifyUser(publicKey string, hash []byte, signature []byte) (bool, error) {
	pub, err := getPubKey(publicKey)
	if err != nil {
		return false, err
	}
	verification, err := pub.Verify(hash, signature)
	if err != nil {
		return verification, err
	}
	return verification, nil
}

func getPubKey(publicKey string) (libp2pcrypto.PubKey, error) {
	pubKey, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	return crypto.RestorePubKey(pubKey)
}
