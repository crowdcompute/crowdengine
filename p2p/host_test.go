package p2p

import (
	"testing"

	"github.com/crowdcompute/crowdengine/common"

	"github.com/stretchr/testify/assert"
)

var (
	testHost4 = NewHost(2002, "127.0.0.1", nil)
	testHost5 = NewHost(2002, "127.0.0.1", nil)
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
