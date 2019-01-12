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

// NewImageManagerAPI creates a new image manager RPC API
func NewImageManagerAPI(h *p2p.Host) *ImageManagerAPI {
	return &ImageManagerAPI{
		host: h,
		// TODO: NOT SURE IF THIS IS A GOOD IDEA
		images: make(map[string][]byte),
	}
}

// PushImage is the API call to push an image to the remote peer
func (api *ImageManagerAPI) PushImage(ctx context.Context, nodePID string, imageFilePath string) (string, error) {

	file, err := os.Open(imageFilePath)
	common.CheckErr(err, "[PushImage] Couldn't open file.")
	defer file.Close()

	peerid, err := peer.IDB58Decode(nodePID)
	common.CheckErr(err, "[PushImage] Couldn't IDB58Decode nodePID.")

	fileInfo, err := file.Stat()
	common.CheckErr(err, "[PushImage] Couldn't find stats.")

	// Hash image
	hash := crypto.HashImagePath(imageFilePath)
	hashString := hex.EncodeToString(hash)
	signature := api.images[hashString]
	signatureString := hex.EncodeToString(signature)

	// Starting a new stream to send a file
	api.host.SetConsistentStream(peerid)

	// TODO: all those numbers should go as constants
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	// TODO: Not sure what number to give here. Need to see the range
	filledSignature := fillString(signatureString, 150)
	// TODO: Not sure what number to give here. Need to see the range
	filledHash := fillString(hashString, 100)
	log.Println("Sending filename and filesize!")
	log.Println("fileSize: ", fileSize)
	log.Println("fileName: ", fileName)
	log.Println("filledSignature: ", filledSignature)
	log.Println("filledHash: ", filledHash)

	api.host.UploadChunk([]byte(fileSize))
	api.host.UploadChunk([]byte(fileName))
	api.host.UploadChunk([]byte(filledSignature))
	api.host.UploadChunk([]byte(filledHash))
	sendBuffer := make([]byte, common.FileChunk)
	log.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		api.host.UploadChunk(sendBuffer)
	}
	log.Println("File has been sent, closing connection!")
	return <-api.host.ImageIDchan, nil
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

// RunImage is the API call to run an imageID to the nodePID node
func (api *ImageManagerAPI) RunImage(ctx context.Context, nodePID string, imageID string) string {
	toNodeID, _ := peer.IDB58Decode(nodePID)
	api.host.RunImage(toNodeID, imageID)

	// Check if there are any pending requests to run
	containerID := <-api.host.ContainerID
	log.Println("Result running the job: ", containerID)
	return containerID
}

func (api *ImageManagerAPI) InspectContainer(ctx context.Context, nodePID string, containerID string) (string, error) {
	toNodeID, _ := peer.IDB58Decode(nodePID)
	api.host.CreateSendInspectRequest(toNodeID, containerID)
	log.Println("Result running the job: ")
	return <-api.host.InspectChan, nil
}

// Uploading an image to the current node
// TODO: Upload an image to the supernode instead of passing the file path
func (api *ImageManagerAPI) UploadImage(ctx context.Context, imageFilePath string, privateKey string) (string, error) {
	privByte, _ := hex.DecodeString(privateKey)
	priv, err := crypto.RestorePrivateKey(privByte)

	// Hash image
	// TODO: bytes will be received straight away, not from a path
	hash := crypto.HashImagePath(imageFilePath)
	// content, _ := ioutil.ReadFile(imageFilePath)
	sign, err := priv.Sign(hash)
	common.CheckErr(err, "[UploadImage] Failed to sign image.")
	api.images[hex.EncodeToString(hash)] = sign
	return "", nil
}

// Getting the list of images specific to the publicKey
func (api *ImageManagerAPI) ListImages(ctx context.Context, nodePID string, publicKey string) (string, error) {
	toNodeID, _ := peer.IDB58Decode(nodePID)
	api.host.CreateAndSendListRequest(toNodeID, publicKey)
	return <-api.host.ListChan, nil
}
