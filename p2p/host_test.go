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
	"github.com/crowdcompute/crowdengine/common"

	"github.com/stretchr/testify/assert"
)

var (
	testHost4, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10209, ListenAddress: "127.0.0.1"},
	})
	testHost5, _ = NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10210, ListenAddress: "127.0.0.1"},
	})
)

func TestConnectWithNodes(t *testing.T) {
	testHost4.ConnectWithNodes([]string{testHost5.FullAddr})
	// Should have itself and the peer 5
	noOfPeers := testHost4.P2PHost.Peerstore().Peers().Len()
	assert.True(t, noOfPeers == 2)
	// Check that the host connected exists in the peerstore
	existsInPeerstore := common.SliceExists(testHost4.P2PHost.Peerstore().Peers(), testHost5.P2PHost.ID())
	assert.True(t, existsInPeerstore)
}
