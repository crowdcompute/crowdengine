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

package p2p

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	uuid "github.com/satori/go.uuid"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
)

const imageUploadRequest = "/image/uploadreq/0.0.1"
const imageUploadResponse = "/image/uploadresp/0.0.1"

// UploadImageProtocol type
type UploadImageProtocol struct {
	p2pHost     host.Host // local host
	stream      inet.Stream
	ImageIDchan chan string
}

func NewUploadImageProtocol(p2pHost host.Host) *UploadImageProtocol {
	p := &UploadImageProtocol{p2pHost: p2pHost,
		ImageIDchan: make(chan string, 1),
	}
	p2pHost.SetStreamHandler(imageUploadRequest, p.onUploadRequest)
	p2pHost.SetStreamHandler(imageUploadResponse, p.onUploadResponse)
	return p
}

func (p *UploadImageProtocol) SetConsistentStream(hostID peer.ID) bool {
	log.Printf("%s: Uploading image. Sending request to: %s....", p.p2pHost.ID(), hostID)
	stream, err := p.p2pHost.NewStream(context.Background(), hostID, imageUploadRequest)
	p.stream = stream
	common.CheckErr(err, "[SetConsistentStream] Couldn't set a new stream.")

	return true
}

func (p *UploadImageProtocol) UploadChunk(chunk []byte) bool {
	if _, err := p.stream.Write(chunk); err != nil {
		log.Println("Error writting to stream", err)
		return false
	}
	return true
}

// remote peer requests handler
func (p *UploadImageProtocol) onUploadRequest(s inet.Stream) {
	log.Printf("%s: Received upload request from: %s.", p.p2pHost.ID(), s.Conn().RemotePeer())
	defer s.Reset()

	log.Println("Start receiving the file name and file size")

	fileSize, fileName, signature, hash := extractDataFromStream(s)
	filePath := common.ImagesDest + fileName
	createNewFileFromStream(s, filePath, fileSize)

	imageID, err := loadImageToDocker(filePath)
	if errRemove := removeFile(filePath); errRemove != nil {
		p.ImageIDchan <- errRemove.Error()
		return
	}
	if err != nil {
		errmsg := fmt.Sprintf("There was an error loading the image. Error: %s\n", err)
		log.Printf(errmsg)
		// if the error is that file does not exist then no need to do anything with the error
		p.ImageIDchan <- errmsg
		return
	}

	p.storeImgDataToDB(imageID, hash, signature)
	p.createSendResponse(s.Conn().RemotePeer(), imageID)
}

func extractDataFromStream(s inet.Stream) (int64, string, string, string) {
	// TODO: all those numbers should go as constants
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	bufferSignature := make([]byte, 150)
	bufferHash := make([]byte, 100)

	s.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	log.Println(fileSize)

	s.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	log.Println(fileName)

	s.Read(bufferSignature)
	signature := strings.Trim(string(bufferSignature), ":")

	s.Read(bufferHash)
	hash := strings.Trim(string(bufferHash), ":")

	return fileSize, fileName, signature, hash
}

func createNewFileFromStream(s inet.Stream, toFilePath string, fileSize int64) {
	newFile, err := os.Create(toFilePath)
	common.CheckErr(err, "[onUploadRequest] Couldn't create a new file.")
	defer newFile.Close()

	var receivedBytes int64

	for {
		// If the file size is smaller than the chunk size or
		// if it's the final chunk then copy it over and break
		if (fileSize - receivedBytes) < common.FileChunk {
			io.CopyN(newFile, s, (fileSize - receivedBytes))
			s.Read(make([]byte, (receivedBytes+common.FileChunk)-fileSize))
			break
		}
		io.CopyN(newFile, s, common.FileChunk)
		receivedBytes += common.FileChunk
	}
	log.Println("File received completely!")
}

// loadImageToDocker takes a path to an image file and loads it to the docker daemon
func loadImageToDocker(filePath string) (string, error) {
	log.Println("Loading this image: ", filePath)
	loadImageResp, err := manager.GetInstance().LoadImage(filePath)
	if err != nil {
		return "", err
	}
	log.Println(loadImageResp)
	if imgID, exists := getImageID(loadImageResp); exists {
		// Docker image ID is 64 characters
		return imgID[:64], err
	}
	// If no image ID exists, we extract the image ID
	// from listing the specific image using its tag
	imageTag := loadImageResp[2 : len(loadImageResp)-5]
	res := getImageSummaryFromTag(imageTag)
	imgID := strings.Replace(res.ID, "sha256:", "", -1)
	log.Println("Loaded image. Image ID: ", imgID)
	return imgID, err
}

// imageIDExists checks if a docker image ID exists in the loadImageResp.
// Docker image is just after the 'sha256:' prefix
func getImageID(loadImageResp string) (string, bool) {
	r, _ := regexp.Compile("sha256:(.*)")
	matches := r.FindAllStringSubmatch(loadImageResp, -1)
	return matches[0][1], len(matches) != 0
}

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

// removeImageFile removes the imgFilePath file from the machine
func removeFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		errmsg := fmt.Sprintf("There was an error removing the file. Error: %s\n", err)
		log.Printf(errmsg)
		return fmt.Errorf(errmsg)
	}
	return nil
}

// storeNewImageToDB stores the new image's data to our level DB
func (p *UploadImageProtocol) storeImgDataToDB(imageID string, hash string, signature string) {
	image := database.ImageLvlDB{Hash: hash, Signature: signature, CreatedTime: time.Now().Unix()}
	database.GetDB().Model(image).Put([]byte(imageID))
}

// Create and send a response to the toPeer note
func (p *UploadImageProtocol) createSendResponse(toPeer peer.ID, response string) bool {
	resp := &api.UploadImageResponse{UploadImageMsgData: NewUploadImageMsgData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
		ImageID: response}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.UploadImageMsgData.MessageData.Sign = signProtoMsg(resp, key)

	// send the response
	return sendMsg(p.p2pHost, toPeer, resp, protocol.ID(imageUploadResponse))
}

func (p *UploadImageProtocol) onUploadResponse(s inet.Stream) {
	data := &api.UploadImageResponse{}
	decodeProtoMessage(data, s)

	// Authenticate integrity and authenticity of the message
	if valid := authenticateProtoMsg(data, data.UploadImageMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}
	log.Printf("%s: Received upload image response from %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer())

	p.ImageIDchan <- data.ImageID
}
