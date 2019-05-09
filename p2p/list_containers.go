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

	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/log"
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"

	"github.com/docker/docker/api/types"
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

	containers, err := p.ListContainersForUser(data.PubKey)
	if err != nil {
		log.Println("Could not List images. Error : ", err)
		return
	}
	containersBytes, err := json.Marshal(containers)
	if err != nil {
		log.Println(err, "Error marshaling image summaries")
		return
	}
	log.Println("Container summaries:", string(containersBytes))
	p.createSendResponse(s.Conn().RemotePeer(), string(containersBytes))
}

// ListContainersForUser list images for the user with the specific publicKey
func (p *ListContainersProtocol) ListContainersForUser(publicKey string) ([]types.Container, error) {
	containers := make([]types.Container, 0)
	allContainers, err := manager.GetInstance().ListContainers()
	if err != nil {
		return nil, fmt.Errorf("Error listing images. Error: %v", err)
	}

	for _, container := range allContainers {
		hash, signatures, err := getImgDataFromDB(container.ImageID)
		if err != nil {
			if err == database.ErrNotFound {
				log.Println("Continuing... ")
				continue
			}
			return nil, err
		}
		// Verify all signatures for the same image
		for _, signature := range signatures {
			signedBytes, err := hex.DecodeString(signature)
			if err != nil {
				return nil, err
			}
			if ok, err := verifyUser(publicKey, hash, signedBytes); ok && err == nil {
				containers = append(containers, container)
				// TODO: Delete those comments. Only for debugging mode
				// } else if !ok {
				// 	log.Println("Could not verify this user. Signature could not be verified by the Public key...")
			} else if err != nil {
				return nil, err
			}
		}
	}
	return containers, nil
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
