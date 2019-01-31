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

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"

	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/stretchr/testify/assert"
)

var (
	// TestHost2 has testHost1 as a peer, but not the other way around
	testHost1, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10209, ListenAddress: "127.0.0.1"},
	})
	testHost2, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10210, ListenAddress: "127.0.0.1",
			Bootstraper: config.Bootstraper{
				Nodes: []string{testHost1.FullAddr},
			},
		},
	})
	testHost3, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10211, ListenAddress: "127.0.0.1"},
	})
)

func TestSignAuthenticate(t *testing.T) {
	req := discoveryRequestMsg(testHost1.P2PHost)

	key := testHost1.P2PHost.Peerstore().PrivKey(testHost1.P2PHost.ID())
	req.DiscoveryMsgData.MessageData.Sign = signProtoMsg(req, key)
	valid := authenticateProtoMsg(req, req.DiscoveryMsgData.MessageData)
	assert.True(t, valid)
}

// TestHost2 sends a message to testHost1
// TestHost2 has testHost1 as a peer
func TestSendMsgFromConnectedPeers(t *testing.T) {
	req := discoveryRequestMsg(testHost2.P2PHost)

	ok := sendMsg(testHost2.P2PHost, testHost1.P2PHost.ID(), req, protocol.ID(discoveryRequest))
	assert.True(t, ok)
}

// TestHost1 sends a message to testHost2
// TestHost1 doesn't have testHost2 as a peer
func TestSendMsgFromUnconnectedPeers(t *testing.T) {
	req := discoveryRequestMsg(testHost3.P2PHost)

	ok := sendMsg(testHost3.P2PHost, testHost1.P2PHost.ID(), req, protocol.ID(discoveryRequest))
	assert.False(t, ok)
}
