package node

import (
	"context"
	"fmt"

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
	fmt.Println("Lenght of host: ", len(api.host.P2PHost.Addrs()))
	for index := 0; index < len(api.host.P2PHost.Addrs()); index++ {
		fmt.Println("", api.host.P2PHost.Addrs()[index])
	}

	pid := peer.IDB58Encode(api.host.P2PHost.ID())
	initialRequest := api.host.InitNodeDiscoveryReq(numberOfNodes, pid)
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
