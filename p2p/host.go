package p2p

import (
	"context"
	"fmt"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/p2p/protocols"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	ps "github.com/libp2p/go-libp2p-peerstore"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
)

type Host struct {
	P2PHost host.Host
	dht     *dht.IpfsDHT
	IP      string

	*protocols.JoinSwarmProtocol
	*protocols.TaskProtocol
	*protocols.DiscoveryProtocol
	*protocols.UploadImageProtocol
	*protocols.InspectContainerProtocol
	*protocols.ListImagesProtocol
}

func NewHost(port int, IP string, bootnodes []string) *Host {
	s := &Host{IP: IP}
	s.makeRandomHost(port, IP)

	if len(bootnodes) > 0 {
		s.connectWithNodes(bootnodes)
	}
	fmt.Print("Here is my p2p ID: ")
	fmt.Printf("/ip4/%s/tcp/%d/ipfs/%s\n", IP, port, s.P2PHost.ID().Pretty())

	s.registerProtocols()
	return s
}

// Registering all Protocols
func (h *Host) registerProtocols() {
	// TODO: PATH has to be in a config
	h.JoinSwarmProtocol = protocols.NewJoinSwarmProtocol(h.P2PHost, h.IP)
	h.DiscoveryProtocol = protocols.NewDiscoveryProtocol(h.P2PHost, h.dht)
	h.TaskProtocol = protocols.NewTaskProtocol(h.P2PHost)
	// Registering the Observer that wants to get notified when the task is done.
	h.TaskProtocol.Register(h.DiscoveryProtocol)
	h.UploadImageProtocol = protocols.NewUploadImageProtocol(h.P2PHost)
	h.InspectContainerProtocol = protocols.NewInspectContainerProtocol(h.P2PHost)
	h.ListImagesProtocol = protocols.NewListImagesProtocol(h.P2PHost)
}

// makeRandomHost creates a libp2p host with a randomly generated identity.
func (h *Host) makeRandomHost(port int, IP string) {
	// Ignoring most errors for brevity
	// priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	// priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	priv, _, err := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	common.CheckErr(err, "[makeRandomHost] Can't GenerateKeyPair ")

	// listen, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", IP, port))
	listen, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	host, _ := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(listen),
		libp2p.Identity(priv),
	)

	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
	dstore := dsync.MutexWrap(ds.NewMapDatastore())

	// Make the DHT
	ctx := context.Background()
	h.dht = dht.NewDHT(ctx, host, dstore)

	// Make the routed host
	h.P2PHost = rhost.Wrap(host, h.dht)

	// Bootstrap the host
	err = h.dht.Bootstrap(ctx)
	common.CheckErr(err, "[makeRandomHost] Couldn't bootstrap the host.")
}

// Establishing a libp2p connection to this nodes' bootnodes
func (h *Host) connectWithNodes(nodes []string) {
	fmt.Println("Connecting to my Bootnodes: ")
	for _, nodeAddr := range nodes {
		h.addAddrToPeerstore(nodeAddr)
	}
}

// addAddrToPeerstore parses a peer multiaddress and adds
// it to the given host's peerstore, so it knows how to
// contact it. It returns the peer ID of the remote peer.
func (h *Host) addAddrToPeerstore(addr string) peer.ID {
	// The following code extracts target's the peer ID from the
	// given multiaddress
	ipfsaddr, err := ma.NewMultiaddr(addr)
	common.CheckErr(err, "[addAddrToPeerstore] NewMultiaddr function failed.")

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	common.CheckErr(err, "[addAddrToPeerstore] ValueForProtocol function failed.")

	peerid, err := peer.IDB58Decode(pid)
	common.CheckErr(err, "[addAddrToPeerstore] IDB58Decode function failed.")
	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add
	// it to the peerstore so LibP2P knows how to contact it
	h.P2PHost.Peerstore().AddAddr(peerid, targetAddr, ps.PermanentAddrTTL)

	return peerid
}

func (h *Host) PeerCount() int {
	return 0
}
