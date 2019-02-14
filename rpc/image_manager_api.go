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
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/log"

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

	file, fileSize, fileName, signature, hash, err := api.getFileData(imageHash)
	common.FatalIfErr(err, "Error getting file data")
	defer file.Close()

	pID, err := peer.IDB58Decode(peerID)
	common.FatalIfErr(err, "Error decoding the peerID")
	common.FatalIfErr(api.host.SetConsistentStream(pID), "Error setting a consistent steam with the remote peer")
	api.sendFileMetadata(fileSize, fileName, signature, hash)
	common.FatalIfErr(api.sendFile(file), "Error sending the file to the remote peer")

	return <-api.host.ImageIDchan
}

// getFileData gets the file's handler, size, name, hash and signature
func (api *ImageManagerAPI) getFileData(imageHash string) (*os.File, string, string, string, string, error) {
	img, err := database.GetImageAccountFromDB(imageHash)
	if err != nil {
		return nil, "", "", "", "", fmt.Errorf("Couldn't find the image on the database")
	}
	file, err := os.Open(img.Path)
	if err != nil {
		return nil, "", "", "", "", err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, "", "", "", "", err
	}
	fileSizeFilled := common.FillString(strconv.FormatInt(fileInfo.Size(), 10), common.FileSizeLength)
	fileNameFilled := common.FillString(fileInfo.Name(), common.FileNameLength)
	log.Println("fileSize: ", fileSizeFilled)
	log.Println("fileName: ", fileNameFilled)

	signatureFilled := common.FillString(img.Signature, common.SignatureLength)
	hashFilled := common.FillString(imageHash, common.HashLength)
	log.Println("filledSignature: ", signatureFilled)
	log.Println("filledHash: ", hashFilled)

	return file, fileSizeFilled, fileNameFilled, signatureFilled, hashFilled, err
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
func (api *ImageManagerAPI) RunImage(ctx context.Context, peerID string, imageID string) string {
	toNodeID, _ := peer.IDB58Decode(peerID)
	api.host.RunImage(toNodeID, imageID)

	// Check if there are any pending requests to run
	containerID := <-api.host.ContainerID
	log.Println("Result running the job: ", containerID)
	return containerID
}

// InspectContainer inspects a container containerID from the peer peerID
func (api *ImageManagerAPI) InspectContainer(ctx context.Context, peerID string, containerID string) (string, error) {
	toNodeID, _ := peer.IDB58Decode(peerID)
	log.Println("About to inspect a container...")
	api.host.InitiateInspectRequest(toNodeID, containerID)
	log.Println("Result running the job: ")
	return <-api.host.InspectChan, nil
}

// ListImages gets a list of images from the peer peerID using the user's publicKey
func (api *ImageManagerAPI) ListImages(ctx context.Context, peerID string, publicKey string) (string, error) {
	toNodeID, _ := peer.IDB58Decode(peerID)
	api.host.InitiateListRequest(toNodeID, publicKey)
	return <-api.host.ListChan, nil
}
