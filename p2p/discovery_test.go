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
	"time"

	"github.com/crowdcompute/crowdengine/common/hexutil"
	"github.com/crowdcompute/crowdengine/crypto"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	host "github.com/libp2p/go-libp2p-host"
	"github.com/stretchr/testify/assert"
)

func discRequest(host host.Host) *api.DiscoveryRequest {
	return &api.DiscoveryRequest{DiscoveryMsgData: NewDiscoveryMsgData("1", true, host),
		Message: api.DiscoveryMessage_DiscoveryReq}
}

func TestRequestExpired(t *testing.T) {
	dproto := discoveryProtocol(2000)
	req := discRequest(dproto.p2pHost)
	dproto.setReqExpiryTime(req, 0)
	time.Sleep(time.Second)
	assert.True(t, dproto.requestExpired(req))
}

func TestMessageReceivedAgain(t *testing.T) {
	dproto := discoveryProtocol(2000)
	req := discRequest(dproto.p2pHost)
	req.DiscoveryMsgData.InitHash = hexutil.Encode(crypto.GetProtoHash(req))
	// Put this request in the received messages list
	dproto.receivedMsg[req.DiscoveryMsgData.InitHash] = 1000
	assert.True(t, dproto.checkMsgReceived(req))
}

// func TestPendingRequests(t *testing.T) {
// 	testHost1, dht1 := host.MakeRandomHost(2000, "127.0.0.1")
// 	testHost2, dht2 := host.MakeRandomHost(2001, "127.0.0.1")
// 	dproto1 := NewDiscoveryProtocol(testHost1, dht1)
// 	dproto2 := NewDiscoveryProtocol(testHost2, dht2)
// 	testHost1.Peerstore().AddAddrs(testHost2.ID(), testHost2.Addrs(), ps.PermanentAddrTTL)
// 	testHost2.Peerstore().AddAddrs(testHost1.ID(), testHost1.Addrs(), ps.PermanentAddrTTL)

// 	req := &api.DiscoveryRequest{DiscoveryMsgData: NewDiscoveryMsgData("1", true, testHost1),
// 		Message: api.DiscoveryMessage_DiscoveryReq}
// 	req.DiscoveryMsgData.InitNodeID = peer.IDB58Encode(testHost2.ID())
// 	dproto1.setReqExpiryTime(req, 10)
// 	dproto1.pendingReq = append(dproto1.pendingReq, req)
// 	dproto1.CheckPendingReq()

// 	assert.Equal(t, testHost2.ID(), <-dproto2.AvailableNodeID)
// }
