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
	"bufio"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	"github.com/docker/docker/api/types"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
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

func NewListImagesProtocol(p2pHost host.Host) *ListImagesProtocol {
	p := &ListImagesProtocol{
		p2pHost:  p2pHost,
		ListChan: make(chan string, 1),
	}
	p2pHost.SetStreamHandler(imageListRequest, p.onListRequest)
	p2pHost.SetStreamHandler(imageListResponse, p.onListResponse)
	return p
}

func (p *ListImagesProtocol) CreateAndSendListRequest(toHostID peer.ID, pubKey string) {
	req := &api.ListImagesRequest{ListImagesMsgData: NewListImagesMsgData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		PubKey: pubKey}
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.ListImagesMsgData.MessageData.Sign = signData(req, key)
	sendMsg(p.p2pHost, toHostID, req, protocol.ID(imageListRequest))
}

func (p *ListImagesProtocol) onListRequest(s inet.Stream) {
	data := &api.ListImagesRequest{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onListRequest] Could not decode data.")
	// Authenticate integrity and authenticity of the message
	if valid := authenticateMessage(data, data.ListImagesMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	imgSummaries, err := p.ListImages(data.PubKey)
	common.CheckErr(err, "[onListRequest] Could not List images.")

	imgSummariesBytes, _ := json.Marshal(imgSummaries)

	log.Printf("Image summaries:")
	log.Printf(string(imgSummariesBytes))
	// ***************************************** //
	// Sending the response back to the sender of the msg
	resp := &api.ListImagesResponse{ListImagesMsgData: NewListImagesMsgData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
		ListResult: string(imgSummariesBytes)}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.ListImagesMsgData.MessageData.Sign = signData(resp, key)

	// send the response
	sendMsg(p.p2pHost, s.Conn().RemotePeer(), resp, protocol.ID(imageListResponse))
}

func (p *ListImagesProtocol) ListImages(publicKey string) ([]types.ImageSummary, error) {
	imgSummaries := make([]types.ImageSummary, 0)
	summaries, err := manager.GetInstance().ListImages(types.ImageListOptions{All: true})
	common.CheckErr(err, "[ListImages] Failed to List images")

	pubKey, _ := hex.DecodeString(publicKey)
	pub, err := crypto.RestorePubKey(pubKey)

	for _, img := range summaries {
		image := database.ImageLvlDB{}
		imgID := strings.Replace(img.ID, "sha256:", "", -1)
		i, err := database.GetDB().Model(image).Get([]byte(imgID))

		image, ok := i.(database.ImageLvlDB)
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

func (p *ListImagesProtocol) onListResponse(s inet.Stream) {
	data := &api.ListImagesResponse{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onListResponse] Could not decode data.")

	// Authenticate integrity and authenticity of the message
	if valid := authenticateMessage(data, data.ListImagesMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}
	log.Printf("data.ListResult")
	log.Printf(data.ListResult)

	log.Printf("%s: Received List response from %s. Message id:%s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.ListImagesMsgData.MessageData.Id)
	p.ListChan <- data.ListResult
}