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
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"

	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"

	"github.com/gogo/protobuf/proto"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
)

const (
	expirationCycle = time.Minute
	clientVersion   = "go-p2p-node/0.0.1"
)

func containerRunning(containerID string) bool {
	cjson, err := manager.GetInstance().InspectContainer(containerID)
	if err != nil {
		log.Println("Error inspecting container. ID : \n", containerID)
		return false
	}
	// If at least one is running then state that I am busy
	if cjson.State.Running {
		return true
	}
	return false
}

func signData(data proto.Message, key crypto.PrivKey) []byte {
	signature, err := signProtoMessage(data, key)
	if err != nil {
		log.Println("failed to sign proto message")
		return nil
	}
	return signature
}

// sign an outgoing p2p message payload
func signProtoMessage(message proto.Message, key crypto.PrivKey) ([]byte, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return key.Sign(data)
}

// Authenticate incoming p2p message
// message: a protobufs go data object
// data: common p2p message data
func authenticateMessage(message proto.Message, data *api.MessageData) bool {
	// store a temp ref to signature and remove it from message data
	// sign is a string to allow easy reset to zero-value (empty string)
	sign := data.Sign
	data.Sign = make([]byte, 0)

	// marshall data without the signature to protobufs3 binary format
	bin, err := proto.Marshal(message)
	if err != nil {
		log.Println(err, "failed to marshal pb message")
		return false
	}

	// restore sig in message data (for possible future use)
	data.Sign = sign

	// restore peer id binary format from base58 encoded node id data
	peerID, err := peer.IDB58Decode(data.NodeId)
	if err != nil {
		log.Println(err, "Failed to decode node id from base58")
		return false
	}
	// verify the data was authored by the signing peer identified by the public key
	// and signature included in the message
	return verifyData(bin, []byte(sign), peerID, data.NodePubKey)
}

// Verify incoming p2p message data integrity
// data: data to verify
// signature: author signature provided in the message payload
// peerID: author peer id from the message payload
// pubKeyData: author public key from the message payload
func verifyData(data []byte, signature []byte, peerID peer.ID, pubKeyData []byte) bool {
	key, err := crypto.UnmarshalPublicKey(pubKeyData)
	if err != nil {
		log.Println(err, "Failed to extract key from message key data")
		return false
	}

	// extract node id from the provided public key
	idFromKey, err := peer.IDFromPublicKey(key)

	if err != nil {
		log.Println(err, "Failed to extract peer id from public key")
		return false
	}

	// verify that message author node id matches the provided node public key
	if idFromKey != peerID {
		log.Println(err, "Node id and provided public key mismatch")
		return false
	}

	res, err := key.Verify(data, signature)
	if err != nil {
		log.Println(err, "Error authenticating data")
		return false
	}

	return res
}

func sendMsg(p2pHost host.Host, neighbourID peer.ID, msg proto.Message, protocol protocol.ID) bool {
	s, err := p2pHost.NewStream(context.Background(), neighbourID, protocol)
	if err != nil {
		log.Println(err)
		return false
	}

	ok := sendProtoMessage(msg, s)

	if !ok {
		return false
	}

	return true
}

// helper method - writes a protobuf go data object to a network stream
// data: reference of protobuf go data object to send (not the object itself)
// s: network stream to write the data to
func sendProtoMessage(data proto.Message, s net.Stream) bool {
	writer := bufio.NewWriter(s)
	enc := protobufCodec.Multicodec(nil).Encoder(writer)
	err := enc.Encode(data)
	if err != nil {
		log.Println("[sendProtoMessage] Failed to encode data", err)
		return false
	}
	writer.Flush()
	return true
}

// NewMessageData ...
// helper method - generate message data shared between all node's p2p protocols
// messageId: unique for requests, copied from request for responses
func NewMessageData(messageId string, gossip bool, p2pHost host.Host) *api.MessageData {
	// Add protobufs bin data for message author public key
	// this is useful for authenticating  messages forwarded by a node authored by another node
	nodePubKey, err := p2pHost.Peerstore().PubKey(p2pHost.ID()).Bytes()
	common.CheckErr(err, "[NewMessageData] Failed to get public key for sender from local peer store.")

	return &api.MessageData{ClientVersion: clientVersion,
		NodeId:     peer.IDB58Encode(p2pHost.ID()),
		NodePubKey: nodePubKey,
		Timestamp:  time.Now().Unix(),
		Id:         messageId,
		Gossip:     gossip}
}

// helper method - generate message data shared between all node's p2p protocols
// messageId: unique for requests, copied from request for responses
func NewDiscoveryMsgData(messageId string, gossip bool, p2pHost host.Host) *api.DiscoveryMsgData {
	return &api.DiscoveryMsgData{
		MessageData: NewMessageData(messageId, gossip, p2pHost),
	}
}

func NewRunImageMsgData(messageId string, gossip bool, p2pHost host.Host) *api.RunImageMsgData {
	return &api.RunImageMsgData{
		MessageData: NewMessageData(messageId, gossip, p2pHost),
	}
}

func NewUploadImageMsgData(messageId string, gossip bool, p2pHost host.Host) *api.UploadImageMsgData {
	return &api.UploadImageMsgData{
		MessageData: NewMessageData(messageId, gossip, p2pHost),
	}
}

func NewInspectContMsgData(messageId string, gossip bool, p2pHost host.Host) *api.InspectContMsgData {
	return &api.InspectContMsgData{
		MessageData: NewMessageData(messageId, gossip, p2pHost),
	}
}

func NewListImagesMsgData(messageId string, gossip bool, p2pHost host.Host) *api.ListImagesMsgData {
	return &api.ListImagesMsgData{
		MessageData: NewMessageData(messageId, gossip, p2pHost),
	}
}
