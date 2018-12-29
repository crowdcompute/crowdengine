package protocols

// This file is the communication Protocol for Joining a Docker Swarm network.

import (
	"bufio"
	"fmt"
	"log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protocols/protomsgs"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	uuid "github.com/satori/go.uuid"
)

// pattern: /protocol-name/request-or-response-message/version
const joinReq = "/swarm/joinreq/0.0.1"
const joinResOK = "/swarm/joinrespOK/0.0.1"
const joinReqToken = "/swarm/joinreqtoken/0.0.1"
const joinResJoined = "/swarm/joinresjoined/0.0.1"

// JoinSwarmProtocol type
type JoinSwarmProtocol struct {
	p2pHost      host.Host                   // local host
	requests     map[string]*api.JoinRequest // used to access request data from response handlers
	done         chan bool                   // only for demo purposes to stop main from terminating
	WorkerToken  string
	ManagerToken string
	managerIP    string
}

func NewJoinSwarmProtocol(p2pHost host.Host, managerIP string) *JoinSwarmProtocol {
	p := &JoinSwarmProtocol{
		p2pHost:   p2pHost,
		requests:  make(map[string]*api.JoinRequest),
		managerIP: managerIP,
		done:      make(chan bool, 1),
	}
	p2pHost.SetStreamHandler(joinReq, p.onJoinRequest)
	p2pHost.SetStreamHandler(joinResOK, p.onJoinResponseOK)
	p2pHost.SetStreamHandler(joinReqToken, p.onJoinReqToken)
	p2pHost.SetStreamHandler(joinResJoined, p.onJoinResJoined)
	return p
}

// SendJoinToNeighbours sends a join swarm message to this nodes' bootnodes
func (p *JoinSwarmProtocol) SendJoinToNeighbours(taskReplicas int) {
	fmt.Println("Sending Join to my Bootnodes")
	neighbours := p.p2pHost.Peerstore().Peers()
	// nodesToSwarm := len(neighbours)

	// if nodesToSwarm != taskReplicas {
	// 	log.Print("Nodes for the swarm are different from the task's replicas")
	// 	return
	// }

	for _, nodeAddr := range neighbours {
		p.Join(nodeAddr)
	}
	for i := 0; i < taskReplicas; i++ {
		<-p.done
		log.Print("One node joined just now!")
	}
	log.Print("SWARM READY!")
}

// Join sends a join Request to a specific <hostID>
// This is the initiation of a Join communication.
func (p *JoinSwarmProtocol) Join(hostID peer.ID) bool {
	log.Printf("%s: Sending join swarm request to: %s....", p.p2pHost.ID(), hostID)

	// create message data
	req := &api.JoinRequest{MessageData: NewMessageData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		Message: api.MessageType_JoinReq}

	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.MessageData.Sign = signData(req, key)

	sendMsg(p.p2pHost, hostID, req, protocol.ID(joinReq))

	// store ref request so response handler has access to it
	p.requests[req.MessageData.Id] = req
	log.Printf("%s: Join swarm to: %s was sent. Message Id: %s, Message: %s", p.p2pHost.ID(), peer.ID(hostID), req.MessageData.Id, req.Message)
	return true
}

// The nodes receives a Join Request, decodes, validates it
// and sends a response if it's ok with joining the Swarm
func (p *JoinSwarmProtocol) onJoinRequest(s net.Stream) {
	// get request data
	data := &api.JoinRequest{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onJoinRequest] Couldn't decode data.")

	log.Printf("%s: Received join swarm request from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.Message)

	valid := authenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	//Check if already part of a swarm
	busy, err := nodePartOfSwarm()
	common.CheckErr(err, "[onJoinRequest] CheckIfNodeBusy couldn't get info for the swarm.")

	fmt.Printf("I am busy: %t", busy)

	// If this node is not busy with another task then it sends a Join OK response to
	// the node that wants to create a Swarm (manager) so that this node can get another message
	// with the join Swarm token.
	if !busy {

		// generate response message
		log.Printf("%s: Sending join swarm response to %s. Message id: %s...", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id)

		resp := &api.JoinResponse{MessageData: NewMessageData(data.MessageData.Id, false, p.p2pHost),
			Message: api.MessageType_JoinResOK}

		// sign the data
		key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
		resp.MessageData.Sign = signData(resp, key)

		// send the response
		if sendMsg(p.p2pHost, s.Conn().RemotePeer(), resp, protocol.ID(joinResOK)) {
			log.Printf("%s: Join swarm response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
		}
	}
}

func nodePartOfSwarm() (bool, error) {
	swarmInfo, err := manager.GetInstance().SwarmInfo()
	fmt.Printf("[checkIfNodeBusy] I have this nodeID: %s \n", swarmInfo.NodeID)
	if swarmInfo.NodeID == "" {
		return false, err
	}
	return true, err
}

// Node receives a Join Ok Response decodes, validates it
// and sends the Join token and address back to the node
func (p *JoinSwarmProtocol) onJoinResponseOK(s net.Stream) {
	data := &api.JoinResponse{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onJoinResponseOK] Could not decode data.")

	valid := authenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// locate request data and remove it if found
	_, ok := p.requests[data.MessageData.Id]
	if ok {
		// remove request from map as we have processed it here
		delete(p.requests, data.MessageData.Id)
	} else {
		log.Println("Failed to locate request data boject for response")
		return
	}

	log.Printf("%s: Received join swarm OK response from %s. Message id:%s. Message: %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id, data.Message)

	if data.Message == api.MessageType_JoinResOK {
		log.Printf("%s: Sending join token to: %s....", p.p2pHost.ID(), data.MessageData.Id)

		// I will need to test this :p.p2pHost.Addrs[0], i used to be node.config.IP
		req := &api.JoinRequest{MessageData: NewMessageData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
			Message: api.MessageType_JoinReqToken, JoinToken: p.WorkerToken, JoinMasterAddr: fmt.Sprintf("%s:%s", p.managerIP, "2377")}

		key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
		req.MessageData.Sign = signData(req, key)

		// send the response
		if sendMsg(p.p2pHost, s.Conn().RemotePeer(), req, protocol.ID(joinReqToken)) {
			log.Printf("%s: Join swarm response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
		}
	}

}

// Node receives a Join Token & address message, decodes, validates it
// and joins the Swarm
func (p *JoinSwarmProtocol) onJoinReqToken(s net.Stream) {

	// get request data
	data := &api.JoinRequest{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onJoinReqToken] Could not decode data.")

	log.Printf("%s: Received join request with Token from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.Message)

	valid := authenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("%s: token is: %s\n", s.Conn().LocalPeer(), data.JoinToken)
	log.Printf("%s: remoteAddrs is: %s\n", s.Conn().LocalPeer(), data.JoinMasterAddr)

	// Join the swarm
	remoteAddrs := []string{data.JoinMasterAddr}
	// I will need to test this :p.p2pHost.Addrs[0], i used to be node.config.IP
	result, err := manager.GetInstance().SwarmJoin(p.managerIP, "", remoteAddrs, data.JoinToken, "0.0.0.0:2377")
	common.CheckErr(err, "[onJoinReqToken] Couldn't join swarm.")

	fmt.Printf("Swarm result: %t\n", result)
	if result {
		log.Printf("%s: Sending joined successfully message to: %s....", p.p2pHost.ID(), data.MessageData.Id)
		resp := &api.JoinResponse{MessageData: NewMessageData(data.MessageData.Id, false, p.p2pHost),
			Message: api.MessageType_JoinResOK}

		key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
		resp.MessageData.Sign = signData(resp, key)

		// send the response
		if sendMsg(p.p2pHost, s.Conn().RemotePeer(), resp, protocol.ID(joinResJoined)) {
			log.Printf("%s: Join swarm response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
		}
	}
}

// Getting a sesponse from a node that they joined the swarm successfully
func (p *JoinSwarmProtocol) onJoinResJoined(s net.Stream) {
	data := &api.JoinResponse{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onJoinResJoined] Couldn't decode data.")

	valid := authenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}
	log.Printf("%s: %s Node just joined the swarm.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	p.done <- true
}
