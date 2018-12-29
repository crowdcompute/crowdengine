package node

import (
	"context"

	"github.com/crowdcompute/crowdengine/p2p"
	peer "github.com/libp2p/go-libp2p-peer"
)

type BootnodesAPI struct {
	host *p2p.Host
}

func NewBootnodesAPI(h *p2p.Host) *BootnodesAPI {
	return &BootnodesAPI{host: h}
}

func (api *BootnodesAPI) SetBootnodes(ctx context.Context, nodes []string) {
	api.host.ConnectWithNodes(nodes)
}

func (api *BootnodesAPI) GetBootnodes(ctx context.Context) (peers []string) {
	for _, v := range []peer.ID(api.host.P2PHost.Peerstore().PeersWithAddrs()) {
		peers = append(peers, v.Pretty())
	}
	return
}
