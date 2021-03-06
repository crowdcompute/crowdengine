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

	"github.com/crowdcompute/crowdengine/p2p"
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
func (api *DiscoveryAPI) Discover(ctx context.Context, numberOfNodes int) (string, error) {
	initialRequest, err := api.host.GetInitialDiscoveryReq()
	if err != nil {
		return "Couldn't get initial discovery request", err
	}
	// TODO: InitHash is a temporary solution. Should be the public key instead.
	api.host.InitializeDiscovery(initialRequest.DiscoveryMsgData.InitHash, numberOfNodes)
	// No neighbour sent me this message, that's why the empty string as a second parameter
	api.host.ForwardMsgToPeers(initialRequest, "")

	return "Discovering nodes...", nil
}
