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

package rpc

import (
	"context"
	"testing"

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/p2p"
	"github.com/stretchr/testify/assert"
)

var (
	testHost1, _ = p2p.NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10209, ListenAddress: "127.0.0.1"},
	})
	testHost2, _ = p2p.NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10210, ListenAddress: "127.0.0.1"},
	})
	testHost3, _ = p2p.NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10211, ListenAddress: "127.0.0.1"},
	})
	api = NewBootnodesAPI(testHost3)
)

// TestSetAndGetBootnodes asserts GetBootnodes returns the same no. of nodes as set by the SetBootnodes method
func TestSetAndGetBootnodes(t *testing.T) {
	ctx := context.Background()
	nodeIDs := []string{testHost1.FullAddr, testHost2.FullAddr}
	api.SetBootnodes(ctx, nodeIDs)
	assert.True(t, len(nodeIDs) == len(api.GetBootnodes(ctx)))
}
