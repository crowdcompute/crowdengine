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
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/crowdcompute/crowdengine/accounts/keystore"
	"github.com/crowdcompute/crowdengine/common/dockerutil"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/log"
	"github.com/crowdcompute/crowdengine/manager"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/p2p"
	peer "github.com/libp2p/go-libp2p-peer"
)

// ImageManagerAPI represents the image manager RPC API
type ImageManagerAPI struct {
	host *p2p.Host
}

// NewImageManagerAPI creates a new RPC service with methods specific for managing docker images & containers.
func NewImageManagerAPI(h *p2p.Host) *ImageManagerAPI {
	return &ImageManagerAPI{
		host: h,
	}
}

// PushImage is the API call to push an image to the peer peerID
func (api *ImageManagerAPI) PushImage(ctx context.Context, peerID string, imageHash string) string {
	log.Println("Pushing an image to the peer : ", peerID)

	file, filepath, fileSize, fileName, signature, hash, err := api.getFileData(imageHash)
	defer removeImage(filepath, imageHash)
	if err != nil {
		msg := fmt.Sprintf("Error getting file data. Error: %s", err)
		log.Println(msg)
		return msg
	}

	pID, err := peer.IDB58Decode(peerID)

	if api.isCurrentNode(pID) {
		// Loading the image to the current node
		log.Println("The Peer ID given is me, I will load the image locally!")
		log.Println(filepath)
		imgID, err := dockerutil.LoadImageToDocker(filepath)
		common.FatalIfErr(err, "Error loading this image to the current node.")
		if err = database.StoreImageToDB(imgID, hash, signature); err != nil {
			log.Error("There was an error storing this image to DB: ", imgID)
		}
		return imgID
	}

	// Sending the image to a remote node
	common.FatalIfErr(err, "Error decoding the peerID")
	common.FatalIfErr(api.host.SetConsistentStream(pID), "Error setting a consistent steam with the remote peer")
	api.sendFileMetadata(fileSize, fileName, signature, hash)
	common.FatalIfErr(api.sendFile(file), "Error sending the file to the remote peer")

	return <-api.host.ImageIDchan
}

// isCurrentNode checks if the given peer ID is the current node
func (api *ImageManagerAPI) isCurrentNode(pID peer.ID) bool {
	return api.host.P2PHost.ID() == pID
}

// Removed the image specified from the disk and the level DB
func removeImage(filepath, hash string) error {
	os.Remove(filepath)
	image := &database.ImageAccount{}
	err := database.GetDB().Model(image).Delete([]byte(hash))
	if err != nil {
		return fmt.Errorf("There was an error deleting the image from lvldb")
	}
	return nil
}

// getFileData gets the file's handler, size, name, hash and signature
func (api *ImageManagerAPI) getFileData(imageHash string) (*os.File, string, string, string, string, string, error) {
	img, err := database.GetImageAccountFromDB(imageHash)
	if err != nil {
		return nil, "", "", "", "", "", fmt.Errorf("Couldn't find the image on the database")
	}
	file, err := os.Open(img.Path)
	if err != nil {
		return nil, "", "", "", "", "", err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, "", "", "", "", "", err
	}
	fileSizeFilled := common.FillString(strconv.FormatInt(fileInfo.Size(), 10), common.FileSizeLength)
	fileNameFilled := common.FillString(fileInfo.Name(), common.FileNameLength)
	log.Println("fileSize: ", fileSizeFilled)
	log.Println("fileName: ", fileNameFilled)

	signatureFilled := common.FillString(img.Signature, common.SignatureLength)
	hashFilled := common.FillString(imageHash, common.HashLength)
	log.Println("filledSignature: ", signatureFilled)
	log.Println("filledHash: ", hashFilled)

	return file, img.Path, fileSizeFilled, fileNameFilled, signatureFilled, hashFilled, err
}

// sendFileMetadata sends the size, name, hash and signature to the peer through the opened stream
func (api *ImageManagerAPI) sendFileMetadata(fileSize, fileName, signature, hash string) error {
	// Start sending the metadata first
	api.host.WriteChunk([]byte(fileSize))
	api.host.WriteChunk([]byte(fileName))
	api.host.WriteChunk([]byte(signature))
	api.host.WriteChunk([]byte(hash))
	return api.host.GetWriterError()
}

// sendFile sends the file's data to the peer through the opened stream
func (api *ImageManagerAPI) sendFile(file *os.File) error {
	sendBuffer := make([]byte, common.FileChunk)
	log.Println("Start sending file!")
	for {
		_, err := file.Read(sendBuffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		api.host.WriteChunk(sendBuffer)
	}
	log.Println("File has been sent, closing connection!")
	return api.host.GetWriterError()
}

// RunImage is the API call to run an imageID to the peerID node
func (api *ImageManagerAPI) RunImage(ctx context.Context, peerID, imageID string) (string, error) {
	pID, _ := peer.IDB58Decode(peerID)
	var containerID string
	var err error
	if api.isCurrentNode(pID) {
		containerID, err = manager.GetInstance().CreateRunContainer(imageID)
	} else {
		api.host.RunImage(pID, imageID)
		// Check if there are any pending requests to run
		containerID = <-api.host.ContainerID
	}
	log.Println("Image is running. Container ID: ", containerID)
	return containerID, err
}

// InspectContainer inspects a container containerID from the peer peerID
func (api *ImageManagerAPI) InspectContainer(ctx context.Context, peerID, containerID string) (string, error) {
	pID, _ := peer.IDB58Decode(peerID)
	var rawInspection string
	var err error
	log.Println("About to inspect a container...")
	if api.isCurrentNode(pID) {
		var rawInspectionBytes []byte
		rawInspectionBytes, err = dockerutil.InspectContainerRaw(containerID)
		rawInspection = string(rawInspectionBytes)
	} else {
		api.host.InitiateInspectRequest(pID, containerID)
		rawInspection = <-api.host.InspectChan
	}
	log.Println("Result of inspecting container: ", rawInspection)
	return rawInspection, err
}

// ListImages gets a list of images from the given <peerID> using the caller's publicKey
func (api *ImageManagerAPI) ListImages(ctx context.Context, peerID string) (string, error) {
	pubBytes, err := getKeyFromContext(ctx)
	if err != nil {
		return "", err
	}
	pID, _ := peer.IDB58Decode(peerID)
	var listImages string
	log.Println("About to list images...")
	if api.isCurrentNode(pID) {
		listImages, err = dockerutil.GetRawImagesForUser(hex.EncodeToString(pubBytes))
	} else {
		api.host.InitiateListImgRequest(pID, hex.EncodeToString(pubBytes))
		listImages = <-api.host.ListImgChan
	}
	return listImages, nil
}

// ListContainers gets a list of containers from a given <peerID> using the caller's publickey
func (api *ImageManagerAPI) ListContainers(ctx context.Context, peerID string) (string, error) {
	pubBytes, err := getKeyFromContext(ctx)
	if err != nil {
		return "", err
	}
	pID, _ := peer.IDB58Decode(peerID)

	var listCont string
	log.Println("About to list containers...")
	if api.isCurrentNode(pID) {
		listCont, err = dockerutil.GetRawContainersForUser(hex.EncodeToString(pubBytes))
	} else {
		api.host.InitiateListContRequest(pID, hex.EncodeToString(pubBytes))
		listCont = <-api.host.ListContChan
	}
	return listCont, err
}

func getKeyFromContext(ctx context.Context) ([]byte, error) {
	key, ok := ctx.Value(common.ContextKeyPair).(*keystore.Key)
	if !ok {
		return nil, fmt.Errorf("There was an error getting the key from the context")
	}
	pubBytes, err := key.KeyPair.Private.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	// Drop first 4 bytes of pub key
	pubBytes = pubBytes[4:]
	return pubBytes, nil
}
