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
	"encoding/hex"
	"testing"
	"time"

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/crypto"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"
	host "github.com/libp2p/go-libp2p-host"
	"github.com/stretchr/testify/assert"
)

var (
	discTestHost3, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10209, ListenAddress: "127.0.0.1"},
	})
	discTestHost4, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10210, ListenAddress: "127.0.0.1"},
	})
	discTestHost5, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10211, ListenAddress: "127.0.0.1"},
	})
	discTestHost6, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10212, ListenAddress: "127.0.0.1"},
	})
	discTestHost7, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10213, ListenAddress: "127.0.0.1"},
	})
	discTestHost8, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10214, ListenAddress: "127.0.0.1"},
	})
	discTestHost9, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10215, ListenAddress: "127.0.0.1"},
	})
	discTestHost10, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10216, ListenAddress: "127.0.0.1"},
	})
)

func discoveryRequestMsg(host host.Host) *api.DiscoveryRequest {
	return &api.DiscoveryRequest{DiscoveryMsgData: NewDiscoveryMsgData("1", true, host),
		Message: api.DiscoveryMessage_DiscoveryReq}
}

// TestSetTTLForDiscReq sets the TTL & expiry for a discovery request and checks if it was set cerrectly
func TestSetTTLForDiscReq(t *testing.T) {
	req := discoveryRequestMsg(discTestHost3.P2PHost)
	now := time.Now()
	ttl := time.Second
	discTestHost3.setTTLForDiscReq(req, ttl)
	// assert.True(t, req.DiscoveryMsgData.TTL == uint32(ttl))
	assert.True(t, req.DiscoveryMsgData.Expiry == uint32(now.Add(ttl).Unix()))
}

// TestMsgExpired tests the requestExpired method which checks if a discovery request got expired
// Setting the TTL to the past
func TestMsgExpired(t *testing.T) {
	req := discoveryRequestMsg(discTestHost4.P2PHost)
	discTestHost4.setTTLForDiscReq(req, -1*time.Second)
	assert.True(t, discTestHost4.requestExpired(req))
}

// TestMsgReceived tests whether a host received a discovery msg
func TestMsgReceived(t *testing.T) {
	req := discoveryRequestMsg(discTestHost6.P2PHost)
	hash, err := crypto.HashProtoMsg(req)
	if err != nil {
		t.Errorf("Failed to HashProtoMsg")
	}
	req.DiscoveryMsgData.InitHash = hex.EncodeToString(hash)
	// mock the reception of this discovery message with the specific hash
	discTestHost5.receivedMsgs[req.DiscoveryMsgData.InitHash] = 0
	assert.True(t, discTestHost5.checkMsgReceived(req))
}

// TestDeleteExpiredMsgs checks that an expired message got deleted
func TestDeleteExpiredMsgs(t *testing.T) {
	req := discoveryRequestMsg(discTestHost7.P2PHost)
	hash, err := crypto.HashProtoMsg(req)
	if err != nil {
		t.Errorf("Failed to HashProtoMsg")
	}
	req.DiscoveryMsgData.InitHash = hex.EncodeToString(hash)
	// Expiry time is in the past
	expiry := uint32(time.Now().Add(-1 * time.Second).Unix())
	discTestHost8.receivedMsgs[req.DiscoveryMsgData.InitHash] = expiry
	assert.True(t, len(discTestHost8.receivedMsgs) == 1)
	discTestHost8.deleteExpiredMsgs()
	assert.True(t, len(discTestHost8.receivedMsgs) == 0)
}

// TestCopyNewDiscoveryRequestHaveDiffSignatures checks that when coping all values of a discovery request
// the signature is different due to different nodes signing it
func TestCopyNewDiscoveryRequestHaveDiffSignatures(t *testing.T) {
	req := discoveryRequestMsg(discTestHost9.P2PHost)
	copiedReq := discTestHost10.copyNewDiscoveryRequest(req)
	reqSignature := string(req.DiscoveryMsgData.MessageData.Sign)
	copiedReqSignature := string(copiedReq.DiscoveryMsgData.MessageData.Sign)
	assert.True(t, reqSignature != copiedReqSignature)
}
