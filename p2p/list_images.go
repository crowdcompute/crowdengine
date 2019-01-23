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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	"github.com/docker/docker/api/types"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	uuid "github.com/satori/go.uuid"
)

const imageListRequest = "/image/ListImgReq/0.0.1"
const imageListResponse = "/image/ListImgResp/0.0.1"

// ListImagesProtocol type
type ListImagesProtocol struct {
	p2pHost  host.Host // local host
	stream   inet.Stream
	ListChan chan string
}

// NewListImagesProtocol sets the protocol's stream handlers and returns a new ListImagesProtocol
func NewListImagesProtocol(p2pHost host.Host) *ListImagesProtocol {
	p := &ListImagesProtocol{
		p2pHost:  p2pHost,
		ListChan: make(chan string, 1),
	}
	p2pHost.SetStreamHandler(imageListRequest, p.onListRequest)
	p2pHost.SetStreamHandler(imageListResponse, p.onListResponse)
	return p
}

// InitiateListRequest sends a list images request to toHostID using the pubKey of the user who initiated it
func (p *ListImagesProtocol) InitiateListRequest(toHostID peer.ID, pubKey string) {
	req := &api.ListImagesRequest{ListImagesMsgData: NewListImagesMsgData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		PubKey: pubKey}
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.ListImagesMsgData.MessageData.Sign = signProtoMsg(req, key)
	sendMsg(p.p2pHost, toHostID, req, protocol.ID(imageListRequest))
}

func (p *ListImagesProtocol) onListRequest(s inet.Stream) {
	data := &api.ListImagesRequest{}
	decodeProtoMessage(data, s)
	// Authenticate integrity and authenticity of the message
	if valid := authenticateProtoMsg(data, data.ListImagesMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	imgSummaries, err := p.listImages(data.PubKey)
	if err != nil {
		log.Printf("Could not List images. Error : ", err)
		return
	}
	imgSummariesBytes, err := json.Marshal(imgSummaries)
	if err != nil {
		log.Println(err, "Error marshaling image summaries")
		return
	}
	log.Println("Image summaries:", string(imgSummariesBytes))
	p.createSendResponse(s.Conn().RemotePeer(), string(imgSummariesBytes))
}

func (p *ListImagesProtocol) listImages(publicKey string) ([]types.ImageSummary, error) {
	imgSummaries := make([]types.ImageSummary, 0)
	summaries, err := manager.GetInstance().ListImages(types.ImageListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("Error listing images. Error: %v", err)
	}
	pubKey, _ := hex.DecodeString(publicKey)
	pub, err := crypto.RestorePubKey(pubKey)

	for _, img := range summaries {
		image := &database.ImageLvlDB{}
		imgID := strings.Replace(img.ID, "sha256:", "", -1)
		i, err := database.GetDB().Model(image).Get([]byte(imgID))

		image, ok := i.(*database.ImageLvlDB)
		if !ok {
			continue
		}
		// If the image was found
		if err == nil {
			hashBytes, err := hex.DecodeString(image.Hash)
			if err != nil {
				log.Println(err, "Error decoding hash")
				return nil, err
			}
			signedBytes, err := hex.DecodeString(image.Signature)
			if err != nil {
				log.Println(err, "Error decoding signature")
				return nil, err
			}
			verification, err := pub.Verify(hashBytes, signedBytes)
			if err != nil {
				log.Println(err, "Error authenticating data")
				return nil, err
			}
			if verification {
				imgSummaries = append(imgSummaries, img)
			}
		}
	}
	return imgSummaries, nil
}

// Create and send a response to the toPeer note
func (p *ListImagesProtocol) createSendResponse(toPeer peer.ID, response string) bool {
	// Sending the response back to the sender of the msg
	resp := &api.ListImagesResponse{ListImagesMsgData: NewListImagesMsgData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
		ListResult: response}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.ListImagesMsgData.MessageData.Sign = signProtoMsg(resp, key)

	// send the response
	return sendMsg(p.p2pHost, toPeer, resp, protocol.ID(imageListResponse))
}

func (p *ListImagesProtocol) onListResponse(s inet.Stream) {
	data := &api.ListImagesResponse{}
	decodeProtoMessage(data, s)

	// Authenticate integrity and authenticity of the message
	if valid := authenticateProtoMsg(data, data.ListImagesMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}
	log.Printf("data.ListResult")
	log.Printf(data.ListResult)

	log.Printf("%s: Received List response from %s. Message id:%s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.ListImagesMsgData.MessageData.Id)
	p.ListChan <- data.ListResult
}
