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

func discoveryRequestMsg(host host.Host) *api.DiscoveryRequest {
	return &api.DiscoveryRequest{DiscoveryMsgData: NewDiscoveryMsgData("1", true, host),
		Message: api.DiscoveryMessage_DiscoveryReq}
}

func TestRequestExpired(t *testing.T) {
	req := discoveryRequestMsg(testHost1.P2PHost)
	testHost1.setReqExpiryTime(req, 0)
	time.Sleep(time.Second)
	assert.True(t, testHost1.requestExpired(req))
}

func TestMessageReceived(t *testing.T) {
	req := discoveryRequestMsg(testHost1.P2PHost)
	req.DiscoveryMsgData.InitHash = hexutil.Encode(crypto.GetProtoHash(req))
	testHost1.receivedMsg[req.DiscoveryMsgData.InitHash] = uint32(time.Now().UnixNano())
	assert.True(t, testHost1.checkMsgReceived(req))
}

// func TestPendingRequests(t *testing.T) {
// 	req := discoveryRequestMsg(testHost1.P2PHost)

// 	req.DiscoveryMsgData.InitNodeID = testHost1.P2PHost.ID().Pretty()

// 	testHost2.setReqExpiryTime(req, 10)
// 	testHost2.pendingReq[req] = struct{}{}
// 	testHost2.onNotify()

// 	// fmt.Println(<-testHost1.NodeID)
// 	assert.Equal(t, testHost2.P2PHost.ID(), <-testHost1.NodeID)
// }
