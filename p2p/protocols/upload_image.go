package protocols

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protocols/protomsgs"
	uuid "github.com/satori/go.uuid"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
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
		fmt.Println("Error writting to stream", err)
		return false
	}
	return true
}

// remote peer requests handler
func (p *UploadImageProtocol) onUploadRequest(s inet.Stream) {
	log.Printf("%s: Received upload request from: %s.", p.p2pHost.ID(), s.Conn().RemotePeer())
	defer s.Reset()

	fmt.Println("Start receiving the file name and file size")
	// TODO: all those numbers should go as constants
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	bufferSignature := make([]byte, 150)
	bufferHash := make([]byte, 100)

	s.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	fmt.Println(fileSize)

	s.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	fmt.Println(fileName)

	s.Read(bufferSignature)
	signature := strings.Trim(string(bufferSignature), ":")

	s.Read(bufferHash)
	hash := strings.Trim(string(bufferHash), ":")

	// TODO: we have to set this as a const
	destFileName := common.ImagesDest + fileName
	newFile, err := os.Create(destFileName)
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
	fmt.Println("Received file completely!")

	imageID, err := loadImage(fileName)
	if err != nil {
		fmt.Printf("There was an error loading the image %s\n", err)
		p.ImageIDchan <- "There was an error loading the image"
		return
	}

	p.storeNewImageToDB(imageID, hash, signature)

	//Sending the response
	resp := &api.UploadImageResponse{UploadImageMsgData: NewUploadImageMsgData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
		ImageID: imageID}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.UploadImageMsgData.MessageData.Sign = signData(resp, key)

	// send the response
	sendMsg(p.p2pHost, s.Conn().RemotePeer(), resp, protocol.ID(imageUploadResponse))
}

func loadImage(filename string) (string, error) {
	fmt.Println("Loading this image: ", filename)
	response, err := manager.GetInstance().LoadImage(filename)
	r, _ := regexp.Compile("sha256:(.*)")
	matches := r.FindAllStringSubmatch(response, -1)
	// No image ID exists, so it's the name:tag
	if len(matches) == 0 {
		imageNameTag := response[2 : len(response)-5]
		fmt.Println(imageNameTag)

		fargs := filters.NewArgs()
		fargs.Add("reference", imageNameTag)

		options := types.ImageListOptions{
			Filters: fargs,
		}

		res, err := manager.GetInstance().ListImages(options)
		if err != nil {
			fmt.Println("error: ", err)
		}
		imgID := strings.Replace(res[0].ID, "sha256:", "", -1)
		log.Println("Loaded image. Image ID: ")
		log.Println(imgID)
		return imgID, err
	} else {
		log.Println("Loaded image. Image ID: ")
		log.Println(matches[0][1][:64])
		return matches[0][1][:64], err
	}
}

func (p *UploadImageProtocol) storeNewImageToDB(imageID string, hash string, signature string) {
	image := database.ImageLvlDB{Hash: hash, Signature: signature, CreatedTime: time.Now().Unix()}
	database.GetDB().Model(image).Put([]byte(imageID))
}

func (p *UploadImageProtocol) onUploadResponse(s inet.Stream) {
	data := &api.UploadImageResponse{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onUploadResponse] Couldn't decode data.")

	// Authenticate integrity and authenticity of the message
	if valid := authenticateMessage(data, data.UploadImageMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}
	log.Printf("%s: Received upload image response from %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer())

	p.ImageIDchan <- data.ImageID
}
