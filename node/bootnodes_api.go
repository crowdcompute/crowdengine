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

package node

import (
	"context"

	"github.com/crowdcompute/crowdengine/p2p"
	peer "github.com/libp2p/go-libp2p-peer"
)

type BootnodesAPI struct {
	host *p2p.Host
}

// NewBootnodesAPI creates a new bootnode API
func NewBootnodesAPI(h *p2p.Host) *BootnodesAPI {
	return &BootnodesAPI{host: h}
}

// SetBootnodes connects the current node with the given nodes
func (api *BootnodesAPI) SetBootnodes(ctx context.Context, nodes []string) {
	api.host.ConnectWithNodes(nodes)
}

// GetBootnodes gets the current nodes connected to the current node
func (api *BootnodesAPI) GetBootnodes(ctx context.Context) (peers []string) {
	for _, v := range []peer.ID(api.host.P2PHost.Peerstore().PeersWithAddrs()) {
		peers = append(peers, v.Pretty())
	}
	return
}
