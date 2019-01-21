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
	"io"
	"os"
	"strconv"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/crowdcompute/crowdengine/p2p"
	peer "github.com/libp2p/go-libp2p-peer"
)

// ImageManagerAPI represents the image manager RPC API
type ImageManagerAPI struct {
	host   *p2p.Host
	images map[string][]byte // image hash -> signature
}

// NewImageManagerAPI creates a new RPC service with methods specific for managing docker images & containers.
func NewImageManagerAPI(h *p2p.Host) *ImageManagerAPI {
	return &ImageManagerAPI{
		host: h,
		// TODO: NOT SURE IF THIS IS A GOOD IDEA
		images: make(map[string][]byte),
	}
}

// UploadImage hashes an image and stores its signature on the hash
// TODO: This will change in the future and authorization + token will be used instead of privaty key passed
func (api *ImageManagerAPI) UploadImage(ctx context.Context, imageFilePath string, privateKey string) error {
	privByte, _ := hex.DecodeString(privateKey)
	priv, err := crypto.RestorePrivateKey(privByte)

	// Hash image
	file, err := os.Open(imageFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	hash := crypto.HashFile(file)
	sign, err := priv.Sign(hash)
	api.images[hex.EncodeToString(hash)] = sign
	return err
}

// PushImage is the API call to push an image to the peer peerID
func (api *ImageManagerAPI) PushImage(ctx context.Context, peerID string, imageFilePath string) string {
	file, fileSize, fileName, signature, hash, err := api.getFileData(imageFilePath)
	common.FatalIfErr(err, "")
	defer file.Close()

	pID, err := peer.IDB58Decode(peerID)
	common.FatalIfErr(err, "Error decoding the peerID")
	common.FatalIfErr(api.host.SetConsistentStream(pID), "Error setting a consistent steam with the remote peer")
	api.sendFileMetadata(fileSize, fileName, signature, hash)
	common.FatalIfErr(api.sendFile(file), "Error sending the file to the remote peer")

	return <-api.host.ImageIDchan
}

// getFileData gets the file's handler, size, name, hash and signature
func (api *ImageManagerAPI) getFileData(imageFilePath string) (*os.File, string, string, string, string, error) {
	file, err := os.Open(imageFilePath)
	if err != nil {
		return nil, "", "", "", "", err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, "", "", "", "", err
	}
	// TODO: all those numbers should go as constants
	fileSizeFilled := common.FillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileNameFilled := common.FillString(fileInfo.Name(), 64)
	log.Println("fileSize: ", fileSizeFilled)
	log.Println("fileName: ", fileNameFilled)

	hash := hex.EncodeToString(crypto.HashFile(file))
	signature := hex.EncodeToString(api.images[hash])
	// TODO: Not sure what number to give here. Need to see the range
	signatureFilled := common.FillString(signature, 150)
	// TODO: Not sure what number to give here. Need to see the range
	hashFilled := common.FillString(hash, 100)
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
