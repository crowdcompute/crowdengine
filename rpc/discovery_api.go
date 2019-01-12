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

	"github.com/crowdcompute/crowdengine/log"
	"github.com/crowdcompute/crowdengine/p2p"
	peer "github.com/libp2p/go-libp2p-peer"
)

// DiscoveryAPI represents the discovery RPC API
type DiscoveryAPI struct {
	host *p2p.Host
}

// NewDiscoveryAPI creates a new DiscoveryAPI
func NewDiscoveryAPI(h *p2p.Host) *DiscoveryAPI {
	return &DiscoveryAPI{host: h}
}

// Discover returns a slice of node IDs in the number of the given numberOfNodes
func (api *DiscoveryAPI) Discover(ctx context.Context, numberOfNodes int) ([]string, error) {
	log.Println("Lenght of host: ", len(api.host.P2PHost.Addrs()))
	for index := 0; index < len(api.host.P2PHost.Addrs()); index++ {
		log.Println("", api.host.P2PHost.Addrs()[index])
	}

	pid := peer.IDB58Encode(api.host.P2PHost.ID())
	log.Println(pid)
	pid2 := api.host.P2PHost.ID().Pretty()
	log.Println(pid2)
	initialRequest := api.host.InitNodeDiscoveryReq(numberOfNodes, pid2)
	// This is the initial forward of this message. No neighbour sent me this message, that's why the empty receivedNeighbour
	api.host.ForwardToNeighbours(initialRequest, "")
	// TODO: this channel has to have a TIMEOUT.
	// TODO: Count the number of nodes that replied
	nodeIDs := make([]string, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		nodeIDs = append(nodeIDs, (<-api.host.NodeID).Pretty())
	}

	return nodeIDs, nil
}
