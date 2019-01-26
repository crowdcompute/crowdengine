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
	"context"
	"encoding/hex"
	"sync"
	"time"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/common"
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

	uuid "github.com/satori/go.uuid"
)

// pattern: /protocol-name/request-or-response-message/version
const discoveryRequest = "/Discovery/discoveryreq/0.0.1"
const discoveryResponse = "/Discovery/discoveryresp/0.0.1"

// DiscoveryProtocol implements Observer interface
type DiscoveryProtocol struct {
	p2pHost            host.Host                          // local host
	dht                *dht.IpfsDHT                       // local host
	receivedMsg        map[string]uint32                  // Store all received msgs, so that we do not re-send them when received again
	pendingReq         map[*api.DiscoveryRequest]struct{} // Store all requests that were unable to be fullfiled at the time the node was busy
	maxPendingReq      uint16                             // The maximum requests the node stores for later process
	NodeIDchan         chan peer.ID                       // a way to return the Node ID to the main form
	DiscoveryCompleted bool                               // manages if the discovery was completed or not
	mu                 sync.Mutex
}

// NewDiscoveryProtocol sets the protocol's stream handlers and returns a new DiscoveryProtocol
func NewDiscoveryProtocol(p2pHost host.Host, dht *dht.IpfsDHT) *DiscoveryProtocol {
	p := &DiscoveryProtocol{
		p2pHost:            p2pHost,
		dht:                dht,
		receivedMsg:        make(map[string]uint32),
		maxPendingReq:      5,
		NodeIDchan:         nil,
		DiscoveryCompleted: false,
	}
	p.pendingReq = map[*api.DiscoveryRequest]struct{}{} //p.maxPendingReq
	// Set the handlers the node will be listening to
	p2pHost.SetStreamHandler(discoveryRequest, p.onDiscoveryRequest)
	p2pHost.SetStreamHandler(discoveryResponse, p.onDiscoveryResponse)
	return p
}

// onNotify checks pendingReq map and sends a response back to the sender
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

// InitializeDiscovery initializes the channel
func (p *DiscoveryProtocol) InitializeDiscovery(numberOfNodes int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.DiscoveryCompleted = false
	p.NodeIDchan = make(chan peer.ID, numberOfNodes)
}

// GetInitialDiscoveryReq initializes the msg request that will be forwarded along the network
// Sets the ID of the node that initiated the discovery request
// Sets the unique hash of the msg request
// Sets the TTL & expiry time of the msg request
func (p *DiscoveryProtocol) GetInitialDiscoveryReq() (*api.DiscoveryRequest, error) {
	req := &api.DiscoveryRequest{DiscoveryMsgData: NewDiscoveryMsgData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		Message: api.DiscoveryMessage_DiscoveryReq}

	// The node initilizing this request is the init Node
	req.DiscoveryMsgData.InitNodeID = p.p2pHost.ID().Pretty()
	hash, err := crypto.HashProtoMsg(req)
	if err != nil {
		return nil, err
	}
	req.DiscoveryMsgData.InitHash = hex.EncodeToString(hash)
	p.setTTLForDiscReq(req, common.TTLmsg)
	return req, err
}

func (p *DiscoveryProtocol) setTTLForDiscReq(req *api.DiscoveryRequest, ttl time.Duration) {
	req.DiscoveryMsgData.TTL = uint32(ttl)
	req.DiscoveryMsgData.Expiry = uint32(time.Now().Add(ttl).Unix())
}

// ForwardMsgToPeers creates a new Discovery request and sends it to its neighbours
func (p *DiscoveryProtocol) ForwardMsgToPeers(request *api.DiscoveryRequest, peerWhoSentMsg peer.ID) {
	req := p.copyNewDiscoveryRequest(request)
	p.sendMsgToPeers(req, peerWhoSentMsg)
}

// copyNewDiscoveryRequest gets a DiscoveryRequest and returns a copy of it,
// but with a new message ID and the current node's signature instead
func (p *DiscoveryProtocol) copyNewDiscoveryRequest(request *api.DiscoveryRequest) *api.DiscoveryRequest {
	req := &api.DiscoveryRequest{DiscoveryMsgData: NewDiscoveryMsgData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		Message: api.DiscoveryMessage_DiscoveryReq}
	req.DiscoveryMsgData.InitNodeID = request.DiscoveryMsgData.InitNodeID
	req.DiscoveryMsgData.TTL = request.DiscoveryMsgData.TTL
	req.DiscoveryMsgData.Expiry = request.DiscoveryMsgData.Expiry
	req.DiscoveryMsgData.InitHash = request.DiscoveryMsgData.InitHash

	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.DiscoveryMsgData.MessageData.Sign = signProtoMsg(req, key)
	return req
}

// sendMsgToPeers sends a discovery request message to all of its peers except itself and the peer peerWhoSentMsg
func (p *DiscoveryProtocol) sendMsgToPeers(req *api.DiscoveryRequest, peerWhoSentMsg peer.ID) {
	// Excluded peers from sending a message
	excludedPeers := map[peer.ID]struct{}{}
	excludedPeers[p.p2pHost.ID()] = struct{}{} // Myself
	excludedPeers[peerWhoSentMsg] = struct{}{} // The peer who sent the message to the current node

	// Send message to peers
	for _, neighbourID := range p.p2pHost.Peerstore().Peers() {
		if _, ok := excludedPeers[neighbourID]; !ok {
			go sendMsg(p.p2pHost, neighbourID, req, protocol.ID(discoveryRequest))
			log.Printf("%s: Discovery message to: %s was sent. Message Id: %s, Message: %s",
				p.p2pHost.ID(), neighbourID, req.DiscoveryMsgData.MessageData.Id, req.Message)
		}
	}
}

// onDiscoveryRequest represents a handler
func (p *DiscoveryProtocol) onDiscoveryRequest(s inet.Stream) {
	// get request data
	data := &api.DiscoveryRequest{}
	decodeProtoMessage(data, s)

	// Log the reception of the message
	log.Printf("%s: Received discovery request from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.Message)

	// If the request's TTL expired or
	// If I received the same message again, I will skip
	if p.requestExpired(data) || p.checkMsgReceived(data) {
		return
	}
	// Storing all the received messages
	p.receivedMsg[data.DiscoveryMsgData.InitHash] = data.DiscoveryMsgData.Expiry

	// Authenticate integrity and authenticity of the message
	if valid := authenticateProtoMsg(data, data.DiscoveryMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	if !p.DiscoveryCompleted {
		// Pass this message to my neighbours
		p.ForwardMsgToPeers(data, s.Conn().RemotePeer())
	}

	// Even if there is possibility that we never send a reply to this Node (because of being busy),
	// we still store it our our Peerstore, because there is high possibility to
	// receive a request again.

	// If the node who sent this message is different than the initPeerID
	// then we need to add the init node to our neighbours before sending the message
	initPeerID, _ := peer.IDB58Decode(data.DiscoveryMsgData.InitNodeID)
	if s.Conn().RemotePeer().String() != initPeerID.String() {
		p.dhtFindAddrAndStore(initPeerID)
	}

	busy, err := NodeBusy()
	common.FatalIfErr(err, "Error on checking if node is busy")
	if busy {
		// Cache the request for a later time
		if uint16(len(p.pendingReq)) < p.maxPendingReq {
			p.pendingReq[data] = struct{}{}
		}
		log.Println("I am busy at the moment. Returning...")
		return
	}

	p.createSendResponse(data)
}

// NodeBusy checks if there is at least one running container then it returns true
// This function removes finished containers as well.
// TODO : the logic of a busy node doesn't have to be this way
func NodeBusy() (bool, error) {
	containers, err := manager.GetInstance().ListContainers()

	for _, container := range containers {
		// If at least one is running then state that I am busy
		if containerRunning(container.ID) {
			return false, err
		}
		// If finished container remove it from docker
		manager.GetInstance().RemoveContainer(container.ID, types.ContainerRemoveOptions{})
		break
	}
	return false, err
}

// createSendResponse creates and sends a response back to the peer who initialized the request
func (p *DiscoveryProtocol) createSendResponse(data *api.DiscoveryRequest) bool {
	// Get the init node ID
	initPeerID, _ := peer.IDB58Decode(data.DiscoveryMsgData.InitNodeID)

	resp := &api.DiscoveryResponse{DiscoveryMsgData: NewDiscoveryMsgData(data.DiscoveryMsgData.MessageData.Id, false, p.p2pHost),
		Message: api.DiscoveryMessage_DiscoveryRes}

	// sign the data
	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.DiscoveryMsgData.MessageData.Sign = signProtoMsg(resp, key)

	log.Printf("%s: Discovery response to: %s was sent. Message Id: %s, Message: %s",
		p.p2pHost.ID(), initPeerID, resp.DiscoveryMsgData.MessageData.Id, resp.Message)

	// send the response
	return sendMsg(p.p2pHost, initPeerID, resp, protocol.ID(discoveryResponse))
}

// requestExpired checks if a request req expired
func (p *DiscoveryProtocol) requestExpired(req *api.DiscoveryRequest) bool {
	now := uint32(time.Now().Unix())
	if req.DiscoveryMsgData.Expiry < now {
		log.Printf("Now: %d, expiry: %d", now, req.DiscoveryMsgData.Expiry)
		log.Println("Message Expired. Dropping message... ")
		return true
	}

	return false
}

// checkMsgReceived checks if request req exists in the receivedMsg slice
func (p *DiscoveryProtocol) checkMsgReceived(req *api.DiscoveryRequest) bool {
	if _, ok := p.receivedMsg[req.DiscoveryMsgData.InitHash]; ok {
		log.Println("Already received this message. Dropping message!")
		return true
	}
	return false
}

// dhtFindAddrAndStore finds a peer from the DHT and stores it to the node's peerstore
func (p *DiscoveryProtocol) dhtFindAddrAndStore(initPeerID peer.ID) error {
	ctx := context.Background()
	initPeerInfo, err := p.dht.FindPeer(ctx, initPeerID)
	if err != nil {
		return err
	}
	log.Println("[DHT] Found this init addresses: ")
	log.Println(initPeerInfo.Addrs)
	log.Println("Adding init node to my neighbours:")
	p.p2pHost.Peerstore().AddAddrs(p.p2pHost.ID(), initPeerInfo.Addrs, ps.PermanentAddrTTL)
	return nil
}

// onDiscoveryResponse is a discovery response stream handler
func (p *DiscoveryProtocol) onDiscoveryResponse(s inet.Stream) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// If discovery request was completed discard all responses
	if p.DiscoveryCompleted {
		return
	}
	data := &api.DiscoveryResponse{}
	decodeProtoMessage(data, s)

	// Authenticate integrity and authenticity of the message
	if valid := authenticateProtoMsg(data, data.DiscoveryMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}

	discoveryPeer := s.Conn().RemotePeer()
	log.Printf("%s: Received discovery response from %s. Message id:%s. Message: %s.", s.Conn().LocalPeer(), discoveryPeer, data.DiscoveryMsgData.MessageData.Id, data.Message)
	p.NodeIDchan <- discoveryPeer
}

// DeleteDiscoveryMsgs checks for expired received messages
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

// deleteExpiredMsgs iterates over all the receivedMsg map
// removing all expired messages
func (p *DiscoveryProtocol) deleteExpiredMsgs() {
	now := uint32(time.Now().Unix())
	for hash, expiry := range p.receivedMsg {
		if expiry < now {
			log.Printf("about to delete this: %s\n", hash)
			delete(p.receivedMsg, hash)
		}
	}
}
