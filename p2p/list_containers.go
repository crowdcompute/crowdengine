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

const containersListRequest = "/container/ListContainersReq/0.0.1"
const containersListResponse = "/container/ListContainersResp/0.0.1"

// ListContainersProtocol type
type ListContainersProtocol struct {
	p2pHost      host.Host // local host
	stream       inet.Stream
	ListContChan chan string
}

// NewListContainersProtocol sets the protocol's stream handlers and returns a new ListContainersProtocol
func NewListContainersProtocol(p2pHost host.Host) *ListContainersProtocol {
	p := &ListContainersProtocol{
		p2pHost:      p2pHost,
		ListContChan: make(chan string, 1),
	}
	p2pHost.SetStreamHandler(containersListRequest, p.onListRequest)
	p2pHost.SetStreamHandler(containersListResponse, p.onListResponse)
	return p
}

// InitiateListContRequest sends a list images request to toHostID using the pubKey of the user who initiated it
func (p *ListContainersProtocol) InitiateListContRequest(toHostID peer.ID, pubKey string) {
	req := &api.ListContainersRequest{ListContainersMsgData: NewListContainersMsgData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		PubKey: pubKey}
	p2pPrivKey := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.ListContainersMsgData.MessageData.Sign = signProtoMsg(req, p2pPrivKey)
	sendMsg(p.p2pHost, toHostID, req, protocol.ID(containersListRequest))
}

func (p *ListContainersProtocol) onListRequest(s inet.Stream) {
	data := &api.ListContainersRequest{}
	decodeProtoMessage(data, s)
	// Authenticate integrity and authenticity of the message
	if valid := authenticateProtoMsg(data, data.ListContainersMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	containersRaw, err := dockerutil.GetRawContainersForUser(data.PubKey)
	if err != nil {
		log.Println("Could not List images. Error : ", err)
		return
	}
	p.createSendResponse(s.Conn().RemotePeer(), containersRaw)
}

// Create and send a response to the toPeer note
func (p *ListContainersProtocol) createSendResponse(toPeer peer.ID, response string) bool {
	// Sending the response back to the sender of the msg
	resp := &api.ListContainersResponse{ListContainersMsgData: NewListContainersMsgData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
		ListResult: response}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.ListContainersMsgData.MessageData.Sign = signProtoMsg(resp, key)

	// send the response
	return sendMsg(p.p2pHost, toPeer, resp, protocol.ID(containersListResponse))
}

func (p *ListContainersProtocol) onListResponse(s inet.Stream) {
	data := &api.ListContainersResponse{}
	decodeProtoMessage(data, s)

	// Authenticate integrity and authenticity of the message
	if valid := authenticateProtoMsg(data, data.ListContainersMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}
	log.Printf("%s: Received List response from %s. Message id:%s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.ListContainersMsgData.MessageData.Id)
	p.ListContChan <- data.ListResult
}
