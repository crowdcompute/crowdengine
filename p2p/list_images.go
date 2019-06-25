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
	"github.com/crowdcompute/crowdengine/log"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	"github.com/crowdcompute/crowdengine/common/dockerutil"

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
	ListImgChan chan string
}

// NewListImagesProtocol sets the protocol's stream handlers and returns a new ListImagesProtocol
func NewListImagesProtocol(p2pHost host.Host) *ListImagesProtocol {
	p := &ListImagesProtocol{
		p2pHost:  p2pHost,
		ListImgChan: make(chan string, 1),
	}
	p2pHost.SetStreamHandler(imageListRequest, p.onListRequest)
	p2pHost.SetStreamHandler(imageListResponse, p.onListResponse)
	return p
}

// InitiateListImgRequest sends a list images request to toHostID using the pubKey of the user who initiated it
func (p *ListImagesProtocol) InitiateListImgRequest(toHostID peer.ID, pubKey string) {
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

	imgSummariesRaw, err := dockerutil.GetRawImagesForUser(data.PubKey)
	if err != nil {
		log.Println("Could not List images. Error : ", err)
		return
	}
	p.createSendResponse(s.Conn().RemotePeer(), imgSummariesRaw)
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
	log.Printf("%s: Received List response from %s. Message id:%s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.ListImagesMsgData.MessageData.Id)
	p.ListImgChan <- data.ListResult
}
