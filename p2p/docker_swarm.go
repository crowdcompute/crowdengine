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

// This file is the communication Protocol for Joining a Docker Swarm network.

import (
	"fmt"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	uuid "github.com/satori/go.uuid"
)

// pattern: /protocol-name/request-or-response-message/version
const joinReq = "/swarm/joinreq/0.0.1"
const joinResOK = "/swarm/joinrespOK/0.0.1"
const joinReqToken = "/swarm/joinreqtoken/0.0.1"
const joinResJoined = "/swarm/joinresjoined/0.0.1"
const leaveReq = "/swarm/leavereq/0.0.1"
const leaveResOK = "/swarm/leaveres/0.0.1"

// SwarmProtocol type
type SwarmProtocol struct {
	p2pHost      host.Host // local host
	done         chan bool // only for demo purposes to stop main from terminating
	WorkerToken  string
	ManagerToken string
	managerIP    string
}

func NewSwarmProtocol(p2pHost host.Host, managerIP string) *SwarmProtocol {
	p := &SwarmProtocol{
		p2pHost:   p2pHost,
		managerIP: managerIP,
		done:      make(chan bool, 1),
	}
	p2pHost.SetStreamHandler(joinReq, p.onJoinRequest)
	p2pHost.SetStreamHandler(joinResOK, p.onJoinResponseOK)
	p2pHost.SetStreamHandler(joinReqToken, p.onJoinReqToken)
	p2pHost.SetStreamHandler(joinResJoined, p.onJoinResJoined)
	p2pHost.SetStreamHandler(leaveReq, p.onLeaveRequest)
	p2pHost.SetStreamHandler(leaveResOK, p.onLeaveResponseOK)
	return p
}

// SendJoinToPeersAndWait sends a join swarm request to it's peers
// And waits until taskReplicas nodes are connected
func (p *SwarmProtocol) SendJoinToPeersAndWait(taskReplicas int) {
	log.Println("Sending Join to my connected peers")
	peers := p.p2pHost.Peerstore().Peers()

	for _, nodeAddr := range peers {
		if p.p2pHost.ID() != nodeAddr {
			p.Join(nodeAddr)
		}
	}
	for i := 0; i < taskReplicas; i++ {
		<-p.done
		log.Print("One node joined just now!")
	}
	log.Print("SWARM READY!")
}

// Join sends a join Request to a specific <hostID>
// This is the initiation of a Join communication.
func (p *SwarmProtocol) Join(hostID peer.ID) bool {
	log.Printf("%s: Sending join swarm request to: %s....", p.p2pHost.ID(), hostID)

	// create message data
	req := &api.JoinRequest{MessageData: NewMessageData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		Message: api.MessageType_JoinReq}

	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.MessageData.Sign = signProtoMsg(req, key)

	sendMsg(p.p2pHost, hostID, req, protocol.ID(joinReq))

	log.Printf("%s: Join swarm to: %s was sent. Message Id: %s, Message: %s", p.p2pHost.ID(), peer.ID(hostID), req.MessageData.Id, req.Message)
	return true
}

// onJoinRequest receives a Join Request, decodes, validates it
// and sends a response if it's ok with joining the Swarm
func (p *SwarmProtocol) onJoinRequest(s net.Stream) {
	log.Printf("%s: Received join swarm request from %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer())

	data := &api.JoinRequest{}
	decodeProtoMessage(data, s)
	if valid := authenticateProtoMsg(data, data.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	//Check if already part of a swarm
	busy, err := nodePartOfSwarm()
	common.CheckErr(err, "[onJoinRequest] CheckIfNodeBusy couldn't get info for the swarm.")

	log.Printf("Am I already part of a swarm: %t", busy)

	// If this node is not busy with another task then it sends a Join OK response to
	// the node that wants to create a Swarm (manager) so that this node can get another message
	// with the join Swarm token.
	if !busy {
		p.createSendResponse(s.Conn().RemotePeer(), api.MessageType_JoinResOK, protocol.ID(joinResOK))
	}
}

func nodePartOfSwarm() (bool, error) {
	swarmInfo, err := manager.GetInstance().SwarmInfo()
	log.Printf("[checkIfNodeBusy] I have this nodeID: %s \n", swarmInfo.NodeID)
	if swarmInfo.NodeID == "" {
		return false, err
	}
	return true, err
}

// Node receives a Join Ok Response decodes, validates it
// and sends the Join token and address back to the node
func (p *SwarmProtocol) onJoinResponseOK(s net.Stream) {
	data := &api.JoinResponse{}
	decodeProtoMessage(data, s)
	valid := authenticateProtoMsg(data, data.MessageData)
	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("%s: Received join swarm OK response from %s. Message id:%s. Message: %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id, data.Message)

	if data.Message == api.MessageType_JoinResOK {
		log.Printf("%s: Sending join token to: %s....", p.p2pHost.ID(), data.MessageData.Id)

		// TODO: User might need some nodes to be Managers and some others Workers. Now all are Workers
		req := &api.JoinRequest{MessageData: NewMessageData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
			Message: api.MessageType_JoinReqToken, JoinToken: p.WorkerToken, JoinMasterAddr: fmt.Sprintf("%s:%s", p.managerIP, "2377")}

		key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
		req.MessageData.Sign = signProtoMsg(req, key)

		// send the response
		if sendMsg(p.p2pHost, s.Conn().RemotePeer(), req, protocol.ID(joinReqToken)) {
			log.Printf("%s: Join swarm response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
		}
	}

}

// Node receives a Join Token & address message, decodes, validates it
// and joins the Swarm
func (p *SwarmProtocol) onJoinReqToken(s net.Stream) {
	log.Printf("%s: Received join request with Token from %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer())

	data := &api.JoinRequest{}
	decodeProtoMessage(data, s)
	valid := authenticateProtoMsg(data, data.MessageData)
	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("%s: token is: %s\n", s.Conn().LocalPeer(), data.JoinToken)
	log.Printf("%s: remoteAddrs is: %s\n", s.Conn().LocalPeer(), data.JoinMasterAddr)

	// Join the swarm
	remoteAddrs := []string{data.JoinMasterAddr}
	// I will need to test this :p.p2pHost.Addrs[0], it used to be node.config.IP
	// TODO: port here should go on config
	joinSwarmResult, err := manager.GetInstance().SwarmJoin(p.managerIP, "", remoteAddrs, data.JoinToken, "0.0.0.0:2377")
	common.CheckErr(err, "[onJoinReqToken] Couldn't join swarm.")

	log.Printf("Join Swarm Result: %t\n", joinSwarmResult)
	if joinSwarmResult {
		p.createSendResponse(s.Conn().RemotePeer(), api.MessageType_JoinRes, protocol.ID(joinResJoined))
	}
}

// Create and send a response to the toPeer note
func (p *SwarmProtocol) createSendResponse(toPeer peer.ID, messageType api.MessageType, protoID protocol.ID) bool {
	// generate response message
	log.Printf("%s: Sending swarm response to %s.", p.p2pHost.ID(), toPeer)

	resp := &api.JoinResponse{MessageData: NewMessageData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
		Message: messageType}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.MessageData.Sign = signProtoMsg(resp, key)

	// send the response
	sentOK := sendMsg(p.p2pHost, toPeer, resp, protoID)
	if sentOK {
		log.Printf("%s: Swarm response to %s sent.", p.p2pHost.ID(), toPeer)
	}
	return sentOK
}

// Getting a sesponse from a node that they joined the swarm successfully
func (p *SwarmProtocol) onJoinResJoined(s net.Stream) {
	data := &api.JoinResponse{}
	decodeProtoMessage(data, s)

	valid := authenticateProtoMsg(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}
	log.Printf("%s: %s Node just joined the swarm.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	p.done <- true
}

// SendLeaveToPeersAndWait sends a leave swarm request to it's peers
// And waits until taskReplicas nodes are connected
func (p *SwarmProtocol) SendLeaveToPeersAndWait(taskReplicas int) {
	log.Println("Sending Leave to my connected peers")
	peers := p.p2pHost.Peerstore().Peers()

	for _, nodeAddr := range peers {
		if p.p2pHost.ID() != nodeAddr {
			p.Leave(nodeAddr)
		}
	}
	for i := 0; i < taskReplicas; i++ {
		<-p.done
		log.Print("One node left just now!")
	}
	log.Print("ALL NODES LEFT THE SWARM!")
}

// Leave sends a leave Request to a specific <hostID>
func (p *SwarmProtocol) Leave(hostID peer.ID) bool {
	log.Printf("%s: Sending leave swarm request to: %s....", p.p2pHost.ID(), hostID)

	// create message data
	req := &api.JoinRequest{MessageData: NewMessageData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		Message: api.MessageType_JoinReq}

	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.MessageData.Sign = signProtoMsg(req, key)

	sendMsg(p.p2pHost, hostID, req, protocol.ID(leaveReq))

	log.Printf("%s: Leave swarm to: %s was sent. Message Id: %s, Message: %s", p.p2pHost.ID(), peer.ID(hostID), req.MessageData.Id, req.Message)
	return true
}

// onLeaveRequest receives a Leave Request, decodes, validates it
// and sends a response if it's ok with leaving the Swarm
func (p *SwarmProtocol) onLeaveRequest(s net.Stream) {
	log.Printf("%s: Received leave swarm request from %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer())

	data := &api.JoinRequest{}
	decodeProtoMessage(data, s)
	if valid := authenticateProtoMsg(data, data.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	//Check if already part of a swarm
	partOfSwarm, err := nodePartOfSwarm()
	common.CheckErr(err, "[onLeaveRequest] CheckIfNodeBusy couldn't get info for the swarm.")

	log.Printf("Am I already part of a swarm: %t", partOfSwarm)

	// If this node is part of a swarm then it can leave the swarm
	if partOfSwarm {
		if _, err := manager.GetInstance().LeaveSwarm(); err != nil {
			return
		}
		p.createSendResponse(s.Conn().RemotePeer(), api.MessageType_JoinResOK, protocol.ID(leaveResOK))
	}
}

// Getting a sesponse from a node that they leaved the swarm successfully
func (p *SwarmProtocol) onLeaveResponseOK(s net.Stream) {
	data := &api.JoinResponse{}
	decodeProtoMessage(data, s)
	valid := authenticateProtoMsg(data, data.MessageData)
	if !valid {
		log.Println("Failed to authenticate message")
		return
	}
	log.Printf("%s: %s Node just left the swarm.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	p.done <- true
}