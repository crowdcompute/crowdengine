package rpc

import (
	"testing"
	"time"

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/p2p"
	"github.com/stretchr/testify/assert"
)

var (
	discTestHost1, _ = p2p.NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10217, ListenAddress: "127.0.0.1"},
	})
	discTestHost2, _ = p2p.NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10218, ListenAddress: "127.0.0.1",
			Bootstraper: config.Bootstraper{
				Nodes: []string{discTestHost1.FullAddr},
			},
		},
	})
	discTestHost3, _ = p2p.NewHost(&config.GlobalConfig{
		P2P: config.P2P{ListenPort: 10219, ListenAddress: "127.0.0.1",
			Bootstraper: config.Bootstraper{
				Nodes: []string{discTestHost2.FullAddr},
			},
		},
	})
)

// TestHostAddedPeersToPeerstore test that discTestHost1 added both discTestHost2 & discTestHost3 to its peerstore
// discTestHost1 added discTestHost2 because it received a message from it and
// added discTestHost3 as part of the DHT process (dhtFindAddrAndStore method)
func TestHostAddedPeersToPeerstore(t *testing.T) {
	req, err := discTestHost3.GetInitialDiscoveryReq()
	if err != nil {
		t.Errorf("Couldn't get initial discovery request")
	}
	// TODO: InitHash is a temporary solution. Should be the public key instead.
	discTestHost3.InitializeDiscovery(req.DiscoveryMsgData.InitHash, 2)
	discTestHost3.ForwardMsgToPeers(req, "")
	// wait for the host to receive the message and process it
	time.Sleep(time.Second * 1)
	assert.True(t, discTestHost1.PeerCount() == 2)
}
