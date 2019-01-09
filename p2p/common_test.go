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
	"testing"

	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	peer "github.com/libp2p/go-libp2p-peer"
	ps "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/stretchr/testify/assert"
)

func discoveryProtocol(port int) *DiscoveryProtocol {
	testHost1 := NewHost(port, "127.0.0.1", nil)
	return NewDiscoveryProtocol(testHost1.P2PHost, testHost1.dht)
}

func TestSignAuthenticate(t *testing.T) {
	testHost1 := NewHost(2000, "127.0.0.1", nil)

	req := &api.JoinRequest{MessageData: NewMessageData("1", true, testHost1.P2PHost),
		Message: api.MessageType_JoinReq}
	key := testHost1.P2PHost.Peerstore().PrivKey(testHost1.P2PHost.ID())
	req.MessageData.Sign = signData(req, key)
	valid := authenticateMessage(req, req.MessageData)
	assert.True(t, valid)
}

func TestSendMsgSuccess(t *testing.T) {
	dproto := discoveryProtocol(2000)
	testHost2 := NewHost(2001, "127.0.0.1", nil)

	dproto.p2pHost.Peerstore().AddAddrs(testHost2.P2PHost.ID(), testHost2.P2PHost.Addrs(), ps.PermanentAddrTTL)

	req := discRequest(dproto.p2pHost)
	req.DiscoveryMsgData.InitNodeID = peer.IDB58Encode(testHost2.P2PHost.ID())
	dproto.setReqExpiryTime(req, 10)

	success := sendMsg(dproto.p2pHost, testHost2.P2PHost.ID(), req, protocol.ID(discoveryRequest))
	assert.True(t, success)
}

func TestSendMsgFail(t *testing.T) {
	dproto := discoveryProtocol(2000)
	testHost2 := NewHost(2001, "127.0.0.1", nil)

	req := discRequest(dproto.p2pHost)
	req.DiscoveryMsgData.InitNodeID = peer.IDB58Encode(testHost2.P2PHost.ID())
	dproto.setReqExpiryTime(req, 10)

	success := sendMsg(dproto.p2pHost, testHost2.P2PHost.ID(), req, protocol.ID(discoveryRequest))
	assert.False(t, success)
}
