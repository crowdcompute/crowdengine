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

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	uuid "github.com/satori/go.uuid"
)

const inspectContainerRequest = "/image/inspectreq/0.0.1"
const inspectContainerResponse = "/image/inspectresp/0.0.1"

// UploadImageProtocol type
type InspectContainerProtocol struct {
	p2pHost     host.Host // local host
	stream      inet.Stream
	InspectChan chan string
}

func NewInspectContainerProtocol(p2pHost host.Host) *InspectContainerProtocol {
	p := &InspectContainerProtocol{p2pHost: p2pHost, InspectChan: make(chan string, 1)}
	p2pHost.SetStreamHandler(inspectContainerRequest, p.onInspectRequest)
	p2pHost.SetStreamHandler(inspectContainerResponse, p.onInspectResponse)
	return p
}

// InitiateInspectRequest sends an inspect request to toHostID for the containerID
func (p *InspectContainerProtocol) InitiateInspectRequest(toHostID peer.ID, containerID string) {
	req := &api.InspectContRequest{InspectContMsgData: NewInspectContMsgData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		ContainerID: containerID}
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.InspectContMsgData.MessageData.Sign = signProtoMsg(req, key)

	sendMsg(p.p2pHost, toHostID, req, protocol.ID(inspectContainerRequest))
}

func (p *InspectContainerProtocol) onInspectRequest(s inet.Stream) {
	log.Println("Received inspect container request...")
	data := &api.InspectContRequest{}
	decodeProtoMessage(data, s)
	// Authenticate integrity and authenticity of the message
	if valid := authenticateProtoMsg(data, data.InspectContMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}
	rawInspection, err := inspectContainerRaw(data.ContainerID)
	common.CheckErr(err, "[onInspectRequest] Could not inspect container.")

	p.createSendResponse(s.Conn().RemotePeer(), string(rawInspection))
}

func inspectContainerRaw(containerID string) ([]byte, error) {
	log.Println("Inspecting this container: ", containerID)
	getSize := true
	inspection, rawData, err := manager.GetInstance().InspectContainerRaw(containerID, getSize)
	log.Printf("Result inspection the container %t\n", inspection.State.Running)
	return rawData, err
}

// Create and send a response to the toPeer node
func (p *InspectContainerProtocol) createSendResponse(toPeer peer.ID, response string) bool {
	// Sending the response back to the sender of the msg
	resp := &api.InspectContResponse{InspectContMsgData: NewInspectContMsgData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
		Inspection: response}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.InspectContMsgData.MessageData.Sign = signProtoMsg(resp, key)

	// send the response
	return sendMsg(p.p2pHost, toPeer, resp, protocol.ID(inspectContainerResponse))
}

func (p *InspectContainerProtocol) onInspectResponse(s inet.Stream) {
	data := &api.InspectContResponse{}
	decodeProtoMessage(data, s)

	// Authenticate integrity and authenticity of the message
	if valid := authenticateProtoMsg(data, data.InspectContMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("%s: Received inspect response from %s. Message id:%s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.InspectContMsgData.MessageData.Id)
	p.InspectChan <- data.Inspection
}
