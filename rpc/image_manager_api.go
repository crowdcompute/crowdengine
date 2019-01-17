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
func (api *ImageManagerAPI) UploadImage(ctx context.Context, imageFilePath string, privateKey string) (string, error) {
	privByte, _ := hex.DecodeString(privateKey)
	priv, err := crypto.RestorePrivateKey(privByte)

	// Hash image
	hash := crypto.HashFile(imageFilePath)
	sign, err := priv.Sign(hash)
	common.CheckErr(err, "[UploadImage] Failed to sign image.")
	api.images[hex.EncodeToString(hash)] = sign
	return "", nil
}

// PushImage is the API call to push an image to the peer peerID

func (api *ImageManagerAPI) PushImage(ctx context.Context, peerID string, imageFilePath string) (string, error) {
	file, fileSize, fileName, signature, hash := api.getFileData(imageFilePath)
	defer file.Close()

	pID, err := peer.IDB58Decode(peerID)
	common.CheckErr(err, "[PushImage] Couldn't IDB58Decode peerID.")
	api.host.SetConsistentStream(pID) // Starting a new stream to with peerID

	api.sendFileMetadata(fileSize, fileName, signature, hash)
	if err := api.sendFile(file); err != nil {
		return "", err
	}
	return <-api.host.ImageIDchan, nil
}

func (api *ImageManagerAPI) getFileData(imageFilePath string) (*os.File, string, string, string, string) {
	file, err := os.Open(imageFilePath)
	common.CheckErr(err, "[getData] Couldn't open file.")
	fileInfo, err := file.Stat()
	common.CheckErr(err, "[getData] Couldn't find stats.")
	// TODO: all those numbers should go as constants
	fileSizeFilled := common.FillString(strconv.FormatInt(fileInfo.Size(), 10), "", 10)
	fileNameFilled := common.FillString(fileInfo.Name(), "", 64)
	log.Println("fileSize: ", fileSizeFilled)
	log.Println("fileName: ", fileNameFilled)

	hash := hex.EncodeToString(crypto.HashFile(imageFilePath))
	signature := hex.EncodeToString(api.images[hash])
	// TODO: Not sure what number to give here. Need to see the range
	signatureFilled := common.FillString(signature, "", 150)
	// TODO: Not sure what number to give here. Need to see the range
	hashFilled := common.FillString(hash, "", 100)
	log.Println("filledSignature: ", signatureFilled)
	log.Println("filledHash: ", hashFilled)

	return file, fileSizeFilled, fileNameFilled, signatureFilled, hashFilled
}

func (api *ImageManagerAPI) sendFileMetadata(fileSize, fileName, signature, hash string) {
	// Start sending the metadata first
	api.host.UploadChunk([]byte(fileSize))
	api.host.UploadChunk([]byte(fileName))
	api.host.UploadChunk([]byte(signature))
	api.host.UploadChunk([]byte(hash))
}

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
		api.host.UploadChunk(sendBuffer)
	}
	log.Println("File has been sent, closing connection!")
	return nil
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

func (api *ImageManagerAPI) InspectContainer(ctx context.Context, peerID string, containerID string) (string, error) {
	toNodeID, _ := peer.IDB58Decode(peerID)
	log.Println("About to inspect a container...")
	api.host.InitiateInspectRequest(toNodeID, containerID)
	log.Println("Result running the job: ")
	return <-api.host.InspectChan, nil
}

// Getting the list of images specific to the publicKey
func (api *ImageManagerAPI) ListImages(ctx context.Context, peerID string, publicKey string) (string, error) {
	toNodeID, _ := peer.IDB58Decode(peerID)
	api.host.InitiateListRequest(toNodeID, publicKey)
	return <-api.host.ListChan, nil
}
