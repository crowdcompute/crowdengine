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
	"context"
	"time"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/common/hexutil"
	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/crowdcompute/crowdengine/manager"
	"github.com/docker/docker/api/types"

	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	ps "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"

	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	uuid "github.com/satori/go.uuid"
)

// pattern: /protocol-name/request-or-response-message/version
const discoveryRequest = "/Discovery/discoveryreq/0.0.1"
const discoveryResponse = "/Discovery/discoveryresp/0.0.1"

// DiscoveryProtocol implements Observer interface
type DiscoveryProtocol struct {
	p2pHost       host.Host                          // local host
	dht           *dht.IpfsDHT                       // local host
	receivedMsg   map[string]uint32                  // Store all received msgs, so that we do not re-send them when received again
	pendingReq    map[*api.DiscoveryRequest]struct{} // Store all requests that were unable to be fullfiled at the time
	maxPendingReq uint16                             // The maximum requests the node stores for later process
	NodeID        chan peer.ID                       // a way to return the Node ID to the main form
}

func NewDiscoveryProtocol(p2pHost host.Host, dht *dht.IpfsDHT) *DiscoveryProtocol {
	log.Println("initializing discovery proto")
	p := &DiscoveryProtocol{
		p2pHost:       p2pHost,
		dht:           dht,
		receivedMsg:   make(map[string]uint32),
		maxPendingReq: 5,
		NodeID:        nil,
	}
	p.pendingReq = map[*api.DiscoveryRequest]struct{}{} //p.maxPendingReq
	// Set the handlers the node will be listening to
	p2pHost.SetStreamHandler(discoveryRequest, p.onDiscoveryRequest)
	p2pHost.SetStreamHandler(discoveryResponse, p.onDiscoveryResponse)
	return p
}

// Start tracking pending requests we received while we were busy
// Sending a reponse back to the requests that I got earlier
func (p *DiscoveryProtocol) onNotify() {
	log.Println(" pending requests: ", p.pendingReq)
	for req := range p.pendingReq {
		if !p.requestExpired(req) {
			log.Println("Request not expired, trying to send response")
			if p.createSendResponse(req) {
				delete(p.pendingReq, req)
			}
		}
	}
}

// Full node creates the initial discovery Request
// Sets the full node's nodeID
// Sets the expiry time of the request
// Sets the unique hash of the request
func (p *DiscoveryProtocol) InitNodeDiscoveryReq(numberOfNodes int, initNodeID string) *api.DiscoveryRequest {
	// initialize the channel
	p.NodeID = make(chan peer.ID, numberOfNodes)
	// create message data
	req := &api.DiscoveryRequest{DiscoveryMsgData: NewDiscoveryMsgData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		Message: api.DiscoveryMessage_DiscoveryReq}

	req.DiscoveryMsgData.InitNodeID = initNodeID
	// This time has to be a const somewhere
	p.setReqExpiryTime(req, 15)
	req.DiscoveryMsgData.InitHash = hexutil.Encode(crypto.GetProtoHash(req))
	return req
}

func (p *DiscoveryProtocol) setReqExpiryTime(req *api.DiscoveryRequest, ttl uint32) {
	now := time.Now()
	req.DiscoveryMsgData.TTL = ttl
	req.DiscoveryMsgData.Expiry = uint32(now.Add(time.Second * time.Duration(req.DiscoveryMsgData.TTL)).Unix())
}

// Creates a new Discovery request and sends it to its neighbours
func (p *DiscoveryProtocol) ForwardToNeighbours(request *api.DiscoveryRequest, receivedNeighbour peer.ID) {
	// We are creating a new request here, so that we can authenticate (through its PubKey and NodeID)
	// the node that sent it.
	req := &api.DiscoveryRequest{DiscoveryMsgData: NewDiscoveryMsgData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		Message: api.DiscoveryMessage_DiscoveryReq}
	req.DiscoveryMsgData.InitNodeID = request.DiscoveryMsgData.InitNodeID
	req.DiscoveryMsgData.TTL = request.DiscoveryMsgData.TTL
	req.DiscoveryMsgData.Expiry = request.DiscoveryMsgData.Expiry
	req.DiscoveryMsgData.InitHash = request.DiscoveryMsgData.InitHash

	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.DiscoveryMsgData.MessageData.Sign = signData(req, key)

	excludePeers := make([]peer.ID, 0)
	excludePeers = append(excludePeers, receivedNeighbour) // The neighbour that sent me the message
	excludePeers = append(excludePeers, p.p2pHost.ID())    // Myself
	for _, neighbourID := range p.p2pHost.Peerstore().Peers() {
		if !common.SliceExists(excludePeers, neighbourID) {
			go sendMsg(p.p2pHost, neighbourID, req, protocol.ID(discoveryRequest))
			log.Printf("%s: Discovery message to: %s was sent. Message Id: %s, Message: %s",
				p.p2pHost.ID(), neighbourID, request.DiscoveryMsgData.MessageData.Id, request.Message)
		}
	}
}

// remote peer requests handler
func (p *DiscoveryProtocol) onDiscoveryRequest(s inet.Stream) {
	// get request data
	data := &api.DiscoveryRequest{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onDiscoveryRequest] Could not decode data.")

	// Log the reception of the message
	log.Printf("%s: Received discovery request from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.Message)

	// If the request's TTL expired or
	// If I received the same message again, I will skip
	if p.requestExpired(data) || p.checkMsgReceived(data) {
		return
	}
	// TODO: We have to periodically empty this table.
	// Storing all the received messages
	p.receivedMsg[data.DiscoveryMsgData.InitHash] = data.DiscoveryMsgData.Expiry

	// Authenticate integrity and authenticity of the message
	if valid := authenticateMessage(data, data.DiscoveryMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// Pass this message to my neighbours
	p.ForwardToNeighbours(data, s.Conn().RemotePeer())

	// Even if there is possibility that we never send a reply to this Node (because of being busy),
	// we still store it our our Peerstore, because there is high possibility to
	// receive a request again.

	// If the node who sent this message is different than the initPeerID
	// then we need to add the init node to our neighbours before sending the message
	initPeerID, _ := peer.IDB58Decode(data.DiscoveryMsgData.InitNodeID)
	if s.Conn().RemotePeer().String() != initPeerID.String() {
		p.dhtFindAddrAndStore(initPeerID)
	}

	// TODO: CHECK HERE IF I AM AVAILABLE FOR A TASK
	//if I am not available for the task, store the request for a later process.
	// Maximum pending jobs
	if NodeBusy() {
		// Cache the request for a later time
		if uint16(len(p.pendingReq)) < p.maxPendingReq {
			p.pendingReq[data] = struct{}{}
		}
		log.Println("I am busy at the moment. Returning...")
		return
	}

	p.createSendResponse(data)
}

// This function removes finished containers as well.
// TODO : this has to be revised
// If there is at least one running container then it returns true
func NodeBusy() bool {
	containers, err := manager.GetInstance().ListContainers()
	common.CheckErr(err, "Error listing containers.")

	// TODO: This logic to be changed...
	for _, container := range containers {
		// If at least one is running then state that I am busy
		if containerRunning(container.ID) {
			return true
		} else { // If finished or whatever delete it
			manager.GetInstance().RemoveContainer(container.ID, types.ContainerRemoveOptions{})
		}
	}
	return false
}

// Create and send a response to the Init note
func (p *DiscoveryProtocol) createSendResponse(data *api.DiscoveryRequest) bool {
	// Get the init node ID
	initPeerID, _ := peer.IDB58Decode(data.DiscoveryMsgData.InitNodeID)

	resp := &api.DiscoveryResponse{DiscoveryMsgData: NewDiscoveryMsgData(data.DiscoveryMsgData.MessageData.Id, false, p.p2pHost),
		Message: api.DiscoveryMessage_DiscoveryRes}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.DiscoveryMsgData.MessageData.Sign = signData(resp, key)

	log.Printf("%s: Discovery response to: %s was sent. Message Id: %s, Message: %s",
		p.p2pHost.ID(), initPeerID, resp.DiscoveryMsgData.MessageData.Id, resp.Message)

	// send the response
	return sendMsg(p.p2pHost, initPeerID, resp, protocol.ID(discoveryResponse))
}

func (p *DiscoveryProtocol) requestExpired(req *api.DiscoveryRequest) bool {
	now := uint32(time.Now().Unix())

	if req.DiscoveryMsgData.Expiry < now {
		log.Printf("Now: %d, expiry: %d", now, req.DiscoveryMsgData.Expiry)
		log.Println("Message Expired. Dropping message... ")
		return true
	}

	return false
}

func (p *DiscoveryProtocol) checkMsgReceived(req *api.DiscoveryRequest) bool {
	if _, ok := p.receivedMsg[req.DiscoveryMsgData.InitHash]; ok {
		log.Println("Already received this message. Dropping message!")
		return true
	}
	return false
}

func (p *DiscoveryProtocol) dhtFindAddrAndStore(initPeerID peer.ID) {
	ctx := context.Background()
	initPeerInfo, err := p.dht.FindPeer(ctx, initPeerID)
	common.CheckErr(err, "[DHT] Error finding this address.")
	log.Println("[DHT] Found this init addresses: ")
	log.Println(initPeerInfo.Addrs)
	log.Println("Adding init node to my neighbours:")
	p.p2pHost.Peerstore().AddAddrs(p.p2pHost.ID(), initPeerInfo.Addrs, ps.PermanentAddrTTL)
}

// Init node gets all responses from its peers
func (p *DiscoveryProtocol) onDiscoveryResponse(s inet.Stream) {
	data := &api.DiscoveryResponse{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onDiscoveryResponse] Could not decode data.")

	// Authenticate integrity and authenticity of the message
	if valid := authenticateMessage(data, data.DiscoveryMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	discoveryPeer := s.Conn().RemotePeer()

	log.Printf("%s: Received discovery response from %s. Message id:%s. Message: %s.", s.Conn().LocalPeer(), discoveryPeer, data.DiscoveryMsgData.MessageData.Id, data.Message)
	p.NodeID <- discoveryPeer
}

// Runs forever or until the node's done
func (p *DiscoveryProtocol) DeleteDiscoveryMsgs(quit <-chan struct{}) {
	// Start a ticker to check for expirations
	ticker := time.NewTicker(expirationCycle)
	defer ticker.Stop()

	// Repeat updates until termination is requested
	for {
		select {
		case <-ticker.C:
			p.deleteExpiredMsgs()

		case <-quit:
			return
		}
	}
}

// expire iterates over all the expiration timestamps, removing all stale
// messages from the map.
func (p *DiscoveryProtocol) deleteExpiredMsgs() {
	now := uint32(time.Now().Unix())
	for hash, expiry := range p.receivedMsg {
		if expiry < now {
			log.Printf("about to delete this: %s\n", hash)
			delete(p.receivedMsg, hash)
		}
	}
}
