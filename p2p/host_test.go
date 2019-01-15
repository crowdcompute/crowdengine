package p2p

import (
	"testing"

	"github.com/crowdcompute/crowdengine/common"

	"github.com/stretchr/testify/assert"
)

var (
	// TestHost2 has testHost1 as a peer, but not the other way around
	testHost4 = NewHost(2002, "127.0.0.1", nil)
	testHost5 = NewHost(2002, "127.0.0.1", nil)
)

func TestConnectWithNodes(t *testing.T) {
	testHost4.ConnectWithNodes([]string{testHost5.FullAddr})

	assert.True(t, testHost4.P2PHost.Peerstore().Peers().Len() == 2)
}

func TestConnectWithNodes2(t *testing.T) {
	exists := common.SliceExists(testHost4.P2PHost.Peerstore().Peers(), testHost5.P2PHost.ID())
	assert.True(t, exists)
}
