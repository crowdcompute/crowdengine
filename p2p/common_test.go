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
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/stretchr/testify/assert"
)

var (
	testHost1 = NewHost(2000, "127.0.0.1", nil)
	testHost2 = NewHost(2001, "127.0.0.1", []string{testHost1.FullAddr})
)

func discoveryProtocol(port int) *DiscoveryProtocol {
	return testHost1.DiscoveryProtocol
}

func TestSignAuthenticate(t *testing.T) {
	req := &api.JoinRequest{MessageData: NewMessageData("1", true, testHost1.P2PHost),
		Message: api.MessageType_JoinReq}
	key := testHost1.P2PHost.Peerstore().PrivKey(testHost1.P2PHost.ID())
	req.MessageData.Sign = signData(req, key)
	valid := authenticateMessage(req, req.MessageData)
	assert.True(t, valid)
}

func TestSendMsgSuccess(t *testing.T) {
	req := discoveryRequestMsg(testHost1.P2PHost)
	req.DiscoveryMsgData.InitNodeID = testHost1.P2PHost.ID().Pretty()
	testHost1.setReqExpiryTime(req, 10)

	ok := sendMsg(testHost1.P2PHost, testHost2.P2PHost.ID(), req, protocol.ID(discoveryRequest))
	assert.True(t, ok)
}

func TestSendMsgFail(t *testing.T) {
	req := discoveryRequestMsg(testHost1.P2PHost)
	// req.DiscoveryMsgData.InitNodeID = peer.IDB58Encode(testHost2.P2PHost.ID())
	testHost1.setReqExpiryTime(req, 10)

	ok := sendMsg(testHost1.P2PHost, testHost2.P2PHost.ID(), req, protocol.ID(discoveryRequest))
	assert.False(t, ok)
}
