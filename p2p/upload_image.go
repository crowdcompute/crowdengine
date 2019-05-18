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
	"strconv"
	"strings"
	"time"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/common/dockerutil"
	"github.com/crowdcompute/crowdengine/database"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	uuid "github.com/satori/go.uuid"

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
	ImageIDchan chan string
	sWriter     *binStreamWriter // libp2p stream writter
}

// binStreamWriter represents the libp2p stream writter
// along with the error occuring when writting to the stream
type binStreamWriter struct {
	s   inet.Stream
	err error
}

// Write writes a chunck to the stream.
func (w *binStreamWriter) write(chunk []byte) {
	if w.err != nil {
		return
	}
	_, w.err = w.s.Write(chunk)
}

// NewUploadImageProtocol sets the protocol's stream handlers and returns a new UploadImageProtocol
func NewUploadImageProtocol(p2pHost host.Host) *UploadImageProtocol {
	p := &UploadImageProtocol{p2pHost: p2pHost,
		ImageIDchan: make(chan string, 1),
		sWriter:     &binStreamWriter{},
	}
	p2pHost.SetStreamHandler(imageUploadRequest, p.onUploadRequest)
	p2pHost.SetStreamHandler(imageUploadResponse, p.onUploadResponse)
	return p
}

// SetConsistentStream sets a new stream to accept data
func (p *UploadImageProtocol) SetConsistentStream(hostID peer.ID) error {
	log.Printf("%s: Uploading image. Sending request to: %s....", p.p2pHost.ID(), hostID)
	stream, err := p.p2pHost.NewStream(context.Background(), hostID, imageUploadRequest)
	p.sWriter.s = stream
	return err
}

// WriteChunk writes the chunk of bytes to the stream
// You can call this function multiple times without worrying about handling the error throughout the uploads
// Call GetWriterError() at the end to get the error
func (p *UploadImageProtocol) WriteChunk(chunk []byte) {
	p.sWriter.write(chunk)
}

// GetWriterError returns the error of the binary Stream Writer
func (p *UploadImageProtocol) GetWriterError() error {
	return p.sWriter.err
}

// remote peer requests handler
func (p *UploadImageProtocol) onUploadRequest(s inet.Stream) {
	log.Printf("%s: Received upload request from: %s.", p.p2pHost.ID(), s.Conn().RemotePeer())
	defer s.Reset()

	log.Println("Start receiving the file name and file size")

	fileSize, fileName, signature, hash := readMetadataFromStream(s)
	filePath := common.ImagesDest + fileName
	err := createFileFromStream(s, filePath, fileSize)
	common.FatalIfErr(err, "Couldn't read from stream when uploading a file")

	imageID, err := dockerutil.LoadImageToDocker(filePath)
	if errRemove := common.RemoveFile(filePath); errRemove != nil {
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

	if err = database.StoreImageToDB(imageID, hash, signature); err != nil {
		log.Error("There was an error storing this image to DB: ", imageID)
	}
	log.Printf("This image %s with this hash {%s} and signature {%s} was stored into the DB \n", imageID, hash, signature)
	p.createSendResponse(s.Conn().RemotePeer(), imageID)
}

// readMetadataFromStream reads the metadata from the stream s
func readMetadataFromStream(s inet.Stream) (int64, string, string, string) {
	bufferFileName := make([]byte, common.FileNameLength)
	bufferFileSize := make([]byte, common.FileSizeLength)
	bufferSignature := make([]byte, common.SignatureLength)
	bufferHash := make([]byte, common.HashLength)

	s.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	s.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	s.Read(bufferSignature)
	signature := strings.Trim(string(bufferSignature), ":")

	s.Read(bufferHash)
	hash := strings.Trim(string(bufferHash), ":")

	return fileSize, fileName, signature, hash
}

// createFileFromStream reads a file's data from the stream s
func createFileFromStream(s inet.Stream, toFilePath string, fileSize int64) error {
	newFile, err := os.Create(toFilePath)
	if err != nil {
		return err
	}
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
	return nil
}

// storeImageToDB stores the new image's data to our level DB
// If image exists it will keep the old signature
func (p *UploadImageProtocol) storeImageToDB(imageID string, hash string, signature string) error {
	signatures := make([]string, 0)
	// In the case the imageID already exists in the database we keep the old signatures and append the new one.
	if image, err := database.GetImageFromDB(imageID); err == nil {
		// TODO: Need to check if hash of the same image ID is going to always be the same
		// hashes = append(hashes, image.Hash)
		signatures = image.Signatures
	}
	signatures = append(signatures, signature)
	image := &database.ImageLvlDB{Hash: hash, Signatures: signatures, CreatedTime: time.Now().Unix()}
	// And because the image ID is the same all the values in DB will be updated with the new ones
	return database.GetDB().Model(image).Put([]byte(imageID))
}

// createSendResponse creates and sends a response to the toPeer note
func (p *UploadImageProtocol) createSendResponse(toPeer peer.ID, response string) bool {
	resp := &api.UploadImageResponse{UploadImageMsgData: NewUploadImageMsgData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
		ImageID: response}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.UploadImageMsgData.MessageData.Sign = signProtoMsg(resp, key)

	// send the response
	return sendMsg(p.p2pHost, toPeer, resp, protocol.ID(imageUploadResponse))
}

// onUploadResponse is an upload response stream handler
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
